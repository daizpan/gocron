package cmd

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
)

type runOptions struct {
	Locale string
}

// NewRunCmd creates a new `gocron run` command
func NewRunCmd() *cobra.Command {
	options := runOptions{}
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run schedule",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSchedule(cmd, options)
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&options.Locale, "timezone", "t", "Asia/Tokyo", "TimeZone")
	return cmd
}

func runSchedule(cmd *cobra.Command, options runOptions) error {
	locale, err := time.LoadLocation(options.Locale)
	if err != nil {
		return err
	}
	c := cron.New(cron.WithLocation(locale))

	for _, j := range Conf.Job {
		j := j
		j.Stdout = cmd.OutOrStdout()
		j.Stderr = cmd.ErrOrStderr()
		if _, err := c.AddJob(j.Schedule, j); err != nil {
			return err
		}
		cmd.Printf("Regist job: %+v\n", j)
	}

	c.Start()
	defer c.Stop()
	log.Printf("Entries: %+v\n", c.Entries())

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch

	return nil
}
