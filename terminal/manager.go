package terminal

import (
	"fmt"
	"sync"
)

type Component interface {
	Render()
	SetPosition(x, y int)
}

type ComponentManager struct {
	components map[string]Component
	layout     *Layout
	mu         sync.Mutex
}

func NewComponentManager() *ComponentManager {
	return &ComponentManager{
		components: make(map[string]Component),
		layout:     NewLayout(),
	}
}

func (cm *ComponentManager) Add(id string, component Component) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.components[id] = component
}

func (cm *ComponentManager) Remove(id string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	delete(cm.components, id)
}

func (cm *ComponentManager) Get(id string) (Component, bool) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	comp, exists := cm.components[id]
	return comp, exists
}

func (cm *ComponentManager) RenderAll() {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	ClearScreenFull()
	for _, comp := range cm.components {
		comp.Render()
	}
}

func (cm *ComponentManager) RenderComponent(id string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if comp, exists := cm.components[id]; exists {
		comp.Render()
	}
}

func (cm *ComponentManager) MoveComponent(id string, x, y int) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if comp, exists := cm.components[id]; exists {
		comp.SetPosition(x, y)
		ClearScreenFull()
		for _, c := range cm.components {
			c.Render()
		}
	}
}

func (cm *ComponentManager) GetLayout() *Layout {
	return cm.layout
}

func (cm *ComponentManager) List() []string {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	ids := make([]string, 0, len(cm.components))
	for id := range cm.components {
		ids = append(ids, id)
	}
	return ids
}

func (cm *ComponentManager) Clear() {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.components = make(map[string]Component)
	ClearScreenFull()
}

// Helper to create a grid layout
func (cm *ComponentManager) CreateGrid(rows, cols int, componentFactory func(row, col int) Component) {
	width := cm.layout.GetWidth() / cols
	height := cm.layout.GetHeight() / rows

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			comp := componentFactory(r, c)
			comp.SetPosition(c*width, r*height)
			id := fmt.Sprintf("grid_%d_%d", r, c)
			cm.Add(id, comp)
		}
	}
}
