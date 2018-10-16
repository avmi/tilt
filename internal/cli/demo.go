package cli

import (
	"context"
	"errors"
	"time"

	"github.com/spf13/cobra"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type demoCmd struct {
}

func (c *demoCmd) register() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "demo",
		Short: "Run the demo script",
	}

	return cmd
}

func (c *demoCmd) run(ctx context.Context, args []string) error {
	analyticsService.Incr("cmd.demo", map[string]string{})
	defer analyticsService.Flush(time.Second)

	demo, err := wireDemo(ctx)
	if err != nil {
		return err
	}

	err = demo.Run(ctx)
	s, ok := status.FromError(err)
	if ok && s.Code() == codes.Unknown {
		return errors.New(s.Message())
	} else if err != nil {
		if err == context.Canceled {
			// Expected case, no need to be loud about it, just exit
			return nil
		}
		return err
	}

	return nil
}
