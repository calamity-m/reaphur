package bot

import (
	"context"
	"log/slog"

	centralproto "github.com/calamity-m/reaphur/proto/v1/central"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/google/uuid"
)

func asyncHandler[E bot.Event](log *slog.Logger, listenerFunc func(e E)) bot.EventListener {
	async := func(e E) {
		log.Debug("Entered async handler")
		go func() {
			defer func() {
				if r := recover(); r != nil {
					log.Error("panic recovered", slog.Any("panic", r))
				}

				log.Debug("Finisheded async handler")
			}()
			listenerFunc(e)
		}()
	}

	log.Debug("Created async handler")
	return bot.NewListenerFunc(async)
}

func handleDMMessageCreate(bot *DiscordBot) func(e *events.DMMessageCreate) {
	return func(e *events.DMMessageCreate) {
		ctx := context.Background()
		bot.logger.InfoContext(ctx, "DM_MESSAGE_CREATE Started")

		// Perform some sanity checks
		if e.Message.Author.Bot {
			bot.logger.InfoContext(ctx, "message was created by a bot", slog.Any("msg", e.Message), slog.Any("author", e.Message.Author))
			return
		}

		if len(e.Message.Content) > 500 {
			bot.logger.ErrorContext(ctx, "someone had more than 500 character message for some... reason", slog.Any("message", e.Message))
			return
		}

		marshalId, err := e.Message.Author.ID.MarshalJSON()
		if err != nil {
			bot.logger.ErrorContext(ctx, "error marshaling snowflake id for user", slog.Any("err", err), slog.Any("id", e.Message.Author.ID))
			return
		}

		// Call the central service and try to run some functions on the user's DM message
		input := &centralproto.CallFnUserInputRequest{
			RequestUserId:    uuid.NewSHA1(uuid.Nil, marshalId).String(),
			RequestUserInput: e.Message.Content,
		}
		output, err := bot.central.CallFnUserInput(ctx, input)
		if err != nil {
			bot.logger.ErrorContext(ctx, "error calling central fn", slog.Any("err", err), slog.Any("id", e.Message.Author.ID))
			return
		}
		bot.logger.InfoContext(ctx, "got output from central", slog.Any("output", output))

		// Create response message to the user's DM channel
		msg, err := e.Client().Rest().CreateMessage(
			e.ChannelID,
			discord.NewMessageCreateBuilder().SetContentf(
				"%s",
				output.ResponseMessage,
			).Build(),
		)
		if err != nil {
			bot.logger.ErrorContext(ctx, "failed to create message", slog.Any("err", err), slog.Any("event", e))
			return
		}

		bot.logger.DebugContext(ctx, "created message successfully", slog.Any("msg", msg))
		bot.logger.InfoContext(ctx, "DM_MESSAGE_CREATE Finished")
	}
}

func handleMessageCreate(d *DiscordBot) func(e *events.MessageCreate) {
	return func(e *events.MessageCreate) {
		ctx := context.Background()
		d.logger.InfoContext(ctx, "MESSAGE_CREATE Started")

		d.logger.InfoContext(ctx, "MESSAGE_CREATE Finished")
	}
}
