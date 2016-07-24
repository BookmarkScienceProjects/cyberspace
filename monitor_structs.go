package main

func NewEnvironment(name string) *Environment {
	return &Environment{
		Name:      name,
		Instances: make(map[string]*Instance),
	}
}

type Environment struct {
	Name      string
	Instances map[string]*Instance
}

func (s *Environment) Add(i *Instance) {
	s.Instances[i.InstanceID] = i
}

func NewStack(name string) *Stack {
	return &Stack{
		Name:         name,
		Environments: make(map[string]*Environment),
		Instances:    make(map[string]*Instance),
	}
}

type Stack struct {
	Name         string
	Environments map[string]*Environment
	Instances    map[string]*Instance
}

func (s *Stack) Add(i *Instance) {
	s.Instances[i.InstanceID] = i
	if i.Environment == "" {
		return
	}
	env, ok := s.Environments[i.Environment]
	if !ok {
		env = NewEnvironment(i.Environment)
		s.Environments[i.Environment] = env
	}
	env.Add(i)
}

func NewCluster(name string) *Cluster {
	return &Cluster{
		Name:      name,
		Stacks:    make(map[string]*Stack),
		Instances: make(map[string]*Instance),
	}
}

type Cluster struct {
	Name      string
	Stacks    map[string]*Stack
	Instances map[string]*Instance
}

func (c *Cluster) Add(i *Instance) {
	c.Instances[i.InstanceID] = i
	if i.Stack == "" {
		return
	}
	stack, ok := c.Stacks[i.Stack]
	if !ok {
		stack = NewStack(i.Stack)
		c.Stacks[i.Stack] = stack
	}
	stack.Add(i)
}
