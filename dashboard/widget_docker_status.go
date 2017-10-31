package dashboard

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	humanize "github.com/dustin/go-humanize"
	docker "github.com/fsouza/go-dockerclient"
	ui "github.com/gizak/termui"
	m2s "github.com/mitchellh/mapstructure"
)

// "github.com/k0kubun/pp"

// DockerStatusWidget is a command launcher
type DockerStatusWidget struct {
	widget     Widget
	gauges     []gaugeModel
	isReady    bool
	disabled   bool
	containers []Container
	client     *docker.Client
}

type gaugeModel struct {
	gauge     *ui.Gauge
	metrics   string
	id        string
	name      string
	container string
	active    bool
}

// NewDockerStatusWidget constructs a New DockerStatusWidget
func NewDockerStatusWidget(wi Widget) (n *DockerStatusWidget, err error) {
	n = new(DockerStatusWidget)
	if err := m2s.Decode(wi.Content, &n.containers); err != nil {
		return nil, err
	}
	endpoint := "unix:///var/run/docker.sock"
	n.client, err = docker.NewClient(endpoint)
	if err != nil {
		return nil, err
	}
	n.widget = wi
	if err = n.buildGauges(); err != nil {
		return nil, err
	}
	return
}

func (n *DockerStatusWidget) buildGauges() (err error) {
	l := len(n.containers)
	for _, v := range n.containers {
		g := ui.NewGauge()
		g.Percent = 0
		g.Width = n.widget.RealWidth
		g.Height = n.widget.RealHeight / l
		g.BorderLabelFg = ui.ColorWhite
		var metrics string
		if v.Metrics == "cpu" {
			metrics = "CPU Usage"
			g.BarColor = ui.ColorGreen
		} else if v.Metrics == "memory" {
			metrics = "Memory Usage"
			g.BarColor = ui.ColorRed
		} else {
			return errors.New(v.Metrics + " is not available for the type of docker_status widget")
		}
		lbl := v.Name + " (" + v.Container + ")" + " - " + metrics
		if n.widget.Title == "" {
			g.BorderLabel = lbl
		} else {
			g.BorderLabel = n.widget.Title + " - " + lbl
		}
		g.Label = "fetching... "
		g.LabelAlign = ui.AlignRight

		var id string
		id, err = n.getContainerIDByName(v.Container)
		if err != nil {
			return
		}
		active := true
		if id == "" {
			active = false
			g.Label = "'" + v.Container + "' is not running "
		}
		n.gauges = append(n.gauges, gaugeModel{
			gauge:     g,
			metrics:   v.Metrics,
			id:        id,
			name:      v.Name,
			container: v.Container,
			active:    active,
		})
	}
	return
}

func (n *DockerStatusWidget) getContainerIDByName(name string) (id string, err error) {
	var containers []docker.APIContainers
	containers, err = n.client.ListContainers(docker.ListContainersOptions{
		All: true,
		Filters: map[string][]string{
			"status": []string{"running"},
		},
	})
	if err != nil {
		return
	}
	for _, c := range containers {
		i := strings.LastIndex(c.Names[0], "/"+name)
		if i > -1 {
			id = c.ID
			break
		}
	}
	return
}

// Activate is the implementation of Widget.Activate
func (n *DockerStatusWidget) Activate() {
}

// Deactivate is the implementation of Widget.Activate
func (n *DockerStatusWidget) Deactivate() {
}

// IsDisabled is the implementation of Widget.IsDisabled
func (n *DockerStatusWidget) IsDisabled() bool {
	return n.disabled
}

// IsReady is the implementation of Widget.IsReady
func (n *DockerStatusWidget) IsReady() bool {
	return n.isReady
}

// GetHighlightenPos is the implementation of Widget.GetHighlightenPos
func (n *DockerStatusWidget) GetHighlightenPos() int {
	return 100
}

// GetWidget is the implementation of widget.Activate
func (n *DockerStatusWidget) GetWidget() []ui.GridBufferer {
	var gauges []ui.GridBufferer
	for _, g := range n.gauges {
		gauges = append(gauges, g.gauge)
	}
	return gauges
}

// Render is the implementation of widget.Render
func (n *DockerStatusWidget) Render() (err error) {
	go func() {
		for _, g := range n.gauges {
			if !g.active {
				continue
			}
			if g.metrics == "cpu" {
				r, err := n.readCPU(g)
				if err != nil {
					panic(err)
				}
				per := strconv.FormatFloat(r, 'f', 2, 64)
				g.gauge.Percent = int(r)
				g.gauge.Label = per + "% "
			} else if g.metrics == "memory" {
				m, l, err := n.readMemory(g.id)
				if err != nil {
					panic(err)
				}
				g.gauge.Percent = m
				lim := humanize.Comma(l / 1000 / 1000)
				g.gauge.Label = "{{percent}}% (" + lim + "MBs) "
			}
		}
	}()
	ui.Render(ui.Body)
	return
}

func (n *DockerStatusWidget) readCPU(g gaugeModel) (cpuPercent float64, err error) {
	stats, err := n.getStats(g.id)
	if err != nil {
		return
	}

	var (
		previousCPU    = stats.PreCPUStats.CPUUsage.TotalUsage
		previousSystem = stats.PreCPUStats.SystemCPUUsage
		// calculate the change for the cpu usage of the container in between readings
		cpuDelta = float64(stats.CPUStats.CPUUsage.TotalUsage) - float64(previousCPU)
		// calculate the change for the entire system between readings
		systemDelta = float64(stats.CPUStats.SystemCPUUsage) - float64(previousSystem)
	)

	if systemDelta > 0.0 && cpuDelta > 0.0 {
		cpuPercent = (cpuDelta / systemDelta) * float64(len(stats.CPUStats.CPUUsage.PercpuUsage)) * 100.0
	}

	return
}

func (n *DockerStatusWidget) readMemory(id string) (usage int, memLimit int64, err error) {
	stats, err := n.getStats(id)
	if err != nil {
		return
	}
	memUsage := int64(stats.MemoryStats.Usage - stats.MemoryStats.Stats.Cache)
	memLimit = int64(stats.MemoryStats.Limit)
	memPercent := float64(memUsage) / float64(memLimit) * 100
	usage = int(memPercent)
	return
}

func (n *DockerStatusWidget) getStats(id string) (s *docker.Stats, err error) {
	errC := make(chan error, 1)
	statsC := make(chan *docker.Stats)

	go func() {
		errC <- n.client.Stats(docker.StatsOptions{
			ID:     id,
			Stats:  statsC,
			Stream: false,
		})
	}()

	s, ok := <-statsC
	if !ok {
		return s, fmt.Errorf("Bad response getting stats for container: %s", id)
	}

	err = <-errC
	if err != nil {
		return s, err
	}
	return s, nil
}

// GetWidth is the implementation of widget.Render
func (n *DockerStatusWidget) GetWidth() int {
	return n.widget.RealWidth
}

// GetHeight is the implementation of widget.Render
func (n *DockerStatusWidget) GetHeight() int {
	return n.widget.RealHeight
}
