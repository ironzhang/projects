package framework

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/ironzhang/matrix/framework/pkg/flags"
	"github.com/ironzhang/matrix/framework/pkg/model"
	"github.com/ironzhang/matrix/jsoncfg"
	"github.com/ironzhang/matrix/tlog"
)

type Options struct {
	ConfigFile       string `json:"config-file" usage:"指定配置文件选项"`
	ConfigExample    string `json:"config-example" usage:"生成配置示例选项"`
	LogConfigFile    string `json:"log-config-file" usage:"指定日志配置文件选项"`
	LogConfigExample string `json:"log-config-example" usage:"生成日志配置示例选项"`
}

type Module interface {
	Name() string
	Init() error
	Fini() error
}

type Runner interface {
	Run(ctx context.Context)
}

type framework struct {
	commandLine *flag.FlagSet
	options     Options
	flags       model.Values
	configs     model.Values
	modules     []Module
}

func (f *framework) parseCommandLine() (err error) {
	if f.commandLine == nil {
		f.commandLine = flag.CommandLine
	}
	if err = f.flags.Register("", &f.options); err != nil {
		return err
	}
	for module, opts := range f.flags.Interfaces() {
		if err = flags.Setup(f.commandLine, opts, module, ""); err != nil {
			return err
		}
	}
	return f.commandLine.Parse(os.Args[1:])
}

func (f *framework) doCommandLine() (err error) {
	var quit bool

	if f.options.ConfigExample != "" {
		if jsoncfg.WriteToFile(f.options.ConfigExample, f.configs.Interfaces()); err != nil {
			return fmt.Errorf("generate config example: %v", err)
		}
		quit = true
	}
	if f.options.LogConfigExample != "" {
		if err = jsoncfg.WriteToFile(f.options.LogConfigExample, tlog.Config{}); err != nil {
			return fmt.Errorf("generate log config example: %v", err)
		}
		quit = true
	}

	if quit {
		os.Exit(0)
	}

	return nil
}

func (f *framework) main() (err error) {
	log := tlog.Std().Sugar()

	// module init
	for _, m := range f.modules {
		if err = m.Init(); err != nil {
			log.Errorw("init", "module", m.Name(), "error", err)
			return fmt.Errorf("init %s: %v", m.Name(), err)
		}
		log.Debugw("init", "module", m.Name())
	}

	// quit signal
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt, os.Kill)
		<-ch
		cancel()
		time.Sleep(10 * time.Second)
		fmt.Fprintf(os.Stderr, "wait 10s, force exit")
		log.Sync()
		os.Exit(-3)
	}()

	// module run
	var wg sync.WaitGroup
	for _, m := range f.modules {
		if r, ok := m.(Runner); ok {
			wg.Add(1)
			go func(r Runner) {
				defer wg.Done()
				r.Run(ctx)
			}(r)
		}
	}
	wg.Wait()

	// module fini
	for i := len(f.modules) - 1; i >= 0; i-- {
		m := f.modules[i]
		if err = m.Fini(); err != nil {
			log.Errorw("fini", "module", m.Name(), "error", err)
			continue
		}
		log.Debugw("fini", "module", m.Name())
	}
	return nil
}

func (f *framework) Main() {
	var err error

	// parse command line
	if err = f.parseCommandLine(); err != nil {
		fmt.Fprintf(os.Stderr, "parse command line: %v\n", err)
		os.Exit(3)
	}

	// do command line
	if err = f.doCommandLine(); err != nil {
		fmt.Fprintf(os.Stderr, "do command line: %v\n", err)
		os.Exit(3)
	}

	// load app config
	if err = loadAppConfig(&f.configs, f.options.ConfigFile); err != nil {
		fmt.Fprintf(os.Stderr, "load app config: %v\n", err)
		os.Exit(3)
	}

	// load log config
	log, err := loadLogConfig(f.options.LogConfigFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load log config: %v\n", err)
		os.Exit(3)
	}
	defer log.Sync()

	// main
	if err = f.main(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(3)
	}
}

func (f *framework) Register(m Module, opts interface{}, cfg interface{}) {
	for _, v := range f.modules {
		if v.Name() == m.Name() {
			panic(fmt.Sprintf("module(%s) duplicate", m.Name()))
		}
	}
	f.modules = append(f.modules, m)

	if opts != nil {
		if err := f.flags.Register(m.Name(), opts); err != nil {
			panic(err)
		}
	}

	if cfg != nil {
		if err := f.configs.Register(m.Name(), cfg); err != nil {
			panic(err)
		}
	}
}

var f = &framework{}

func Main() {
	f.Main()
}

func Register(m Module, opts interface{}, cfg interface{}) {
	f.Register(m, opts, cfg)
}

func SetCommandLine(commandLine *flag.FlagSet) {
	f.commandLine = commandLine
}

func SetOptions(opts Options) {
	f.options = opts
}

func Flags() *model.Values {
	return &f.flags
}

func Configs() *model.Values {
	return &f.configs
}
