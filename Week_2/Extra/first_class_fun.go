// Functional options recap

// To begin, I’ll very quickly recap the functional options pattern

type Config struct{ ... }

func WithReticulatedSplines(c *Config) { ... }

type Terrain struct {
        config Config
}

func NewTerrain(options ...func(*Config)) *Terrain {
        var t Terrain
        for _, option := range options {
                option(&t.config)
        }
        return &t 
}

func main() {
        t := NewTerrain(WithReticulatedSplines)
        // [ simulation intensifies ]
}

 // WithCities adds n cities to the Terrain model
 func WithCities(n int) func(*Config) { ... }

 func main() {        
		 t := NewTerrain(WithCities(9))      

// Functions as first class values

// What’s going on here? Let’s break it down. Fundamentally, evaluating a function returns a value. We have functions that take two numbers and return a number.

package math

func Min(a, b float64) float64

// We have functions that take a slice, and return a pointer to a structure.

package bytes

func NewReader(b []byte) *Reader

// and now we have a function which returns a function.

func WithCities(n int) func(*Config)
interface.Apply

// Another way to think about what is going on here is to try to rewrite the functional option pattern using an interface.

type Option interface {
        Apply(*Config)
}

// Rather than a function type we declare an interface, we’ll call it Option, and give it a single method, Apply which takes a pointer to a Config.

func NewTerrain(options ...Option) *Terrain {
        var config Config
        for _, option := range options {
                option.Apply(&config)
        }
        // ...
}

type splines struct{}

func (s *splines) Apply(c *Config) { ... }

func WithReticulatedSplines() Option {
        return new(splines)
}

To write WithCities using our Option interface we need to do a bit more work.

type cities struct {
        cities int
}

func (c *cities) Apply(c *Config) { ... }

func WithCities(n int) Option {
        return &cities{
                cities: n,
        }
}

func main() {
	t := NewTerrain(WithReticulatedSplines(), WithCities(9))
	// ...
}