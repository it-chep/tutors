package main

import (
	"context"

	"github.com/it-chep/tutors.git/internal"
)

func main() {
	ctx := context.Background()

	internal.New(ctx).Run(ctx)
}
