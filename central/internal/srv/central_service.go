package srv

import (
	"context"
	"log/slog"

	"github.com/calamity-m/reaphur/central/internal/fncall"
	"github.com/calamity-m/reaphur/central/internal/parser"
	"github.com/calamity-m/reaphur/pkg/errs"
	centralproto "github.com/calamity-m/reaphur/proto/v1/central"
)

// Simple RPC
//
// Translates user input and actions some user input in some way that the caller cannot know.
// The response will encode releveant information to respond to the user with, but generally
// you cannot know exactly what actions are taken.
// This rpc will parse user input into some structured format, and then
// link together the services required based on that structured output.
func (s *CentralServiceServer) ActionUserInput(ctx context.Context, r *centralproto.ActionUserInputRequest) (*centralproto.ActionUserInputResponse, error) {
	if err := s.commonServiceValidation(); err != nil {
		return nil, err
	}

	parsedBytes, err := s.parser.ActionStructuredOutput(ctx, parser.CreateGenericStructuredOutputRequest(r.RequestUserInput))
	if err != nil {
		s.logger.ErrorContext(ctx, "encountered error calling parser", slog.Any("err", err))
		return nil, err
	}

	s.logger.DebugContext(ctx, "received bytes from parser", slog.Any("bytes", parsedBytes))

	return nil, errs.ErrNotImplementedYet
}

// Simple RPC
//
// As opposed to the actioning of user input, this endpoint allows for the translator
// to actually call the functions themselves, rather than them being stitched together
// by the implementing rpc service.
func (s *CentralServiceServer) CallFnUserInput(ctx context.Context, r *centralproto.CallFnUserInputRequest) (*centralproto.CallFnUserInputResponse, error) {
	if err := s.commonServiceValidation(); err != nil {
		return nil, err
	}

	out, err := s.fnCaller.EnactUserInput(ctx, fncall.CreateGenericFnCallOutputRequest(r.RequestUserInput), s)
	if err != nil {
		s.logger.ErrorContext(ctx, "encountered error calling fn caller", slog.Any("err", err))
		return nil, err
	}

	s.logger.DebugContext(ctx, "received output from fn caller", slog.Any("out", out))

	return nil, errs.ErrNotImplementedYet
}
