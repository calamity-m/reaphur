package bot

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/calamity-m/reaphur/discord/internal/conf"
	centralproto "github.com/calamity-m/reaphur/proto/v1/central"
	"github.com/disgoorg/disgo"
	disgobot "github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/gateway"
)

type DiscordBot struct {
	logger  *slog.Logger
	disc    disgobot.Client
	central centralproto.CentralServiceClient
}

// Connects and runs the discord bot, blocking until the notify chan is pushed to
func (bot *DiscordBot) Run(notify <-chan os.Signal) error {
	// Sync our commands as the first thing we do.
	if err := bot.SyncGlobalCommands(context.TODO()); err != nil {
		bot.logger.Error("failed to sync global commands", slog.Any("err", err))
		return err
	}

	// Now register any listeners we have.
	bot.disc.AddEventListeners([]disgobot.EventListener{
		asyncHandler(bot.logger, handleDMMessageCreate(bot)),
		asyncHandler(bot.logger, handleMessageCreate(bot)),
	}...)

	// Connect to the gateway and defer closing for if we exit
	if err := bot.disc.OpenGateway(context.TODO()); err != nil {
		bot.logger.Error("failed opening discord gateway", slog.Any("err", err))
		return err
	}
	defer bot.Close(context.Background())

	// Block until we receive a notification to stop
	bot.logger.Info("started listening with discord bot")
	<-notify

	bot.logger.Info("finished running discord bot")
	return nil
}

func (bot *DiscordBot) Close(ctx context.Context) {
	bot.logger.Debug("closing")
	bot.disc.Close(ctx)
}

// Retrieves a current list of registered commands from discord. Note that discord global command propogation
// to the discord client itself, takes a long time. (1 hour+)
func (bot *DiscordBot) FetchGlobalCommands(ctx context.Context) ([]discord.ApplicationCommand, error) {
	cmds, err := bot.disc.Rest().GetGlobalCommands(bot.disc.ApplicationID(), true)

	if err != nil {
		return nil, err
	}

	return cmds, nil
}

// Syncs commands from DiscordCmds with discord as global commands. Any global commands that are registered, but not listed
// in the DiscordCmds var will be deleted.
func (bot *DiscordBot) SyncGlobalCommands(ctx context.Context) error {
	// Relies on the doc from discord:
	// "Commands that do not already exist will count toward daily application command create limits."
	//
	// Due to this, we can just blatanly set our global commands at the start, and then just do simple
	// cleanup after.
	bot.logger.InfoContext(ctx, "attempting to sync commands")

	// First, create all of our commands
	set, err := bot.disc.Rest().SetGlobalCommands(
		bot.disc.ApplicationID(),
		BotCommands,
	)
	if err != nil {
		bot.logger.ErrorContext(ctx, "failed setting global commands", slog.Any("commands", BotCommands))
		return fmt.Errorf("failed to create global commands: %w", err)
	}

	bot.logger.InfoContext(ctx, "commands have been set", slog.Any("commands", set))

	return nil
}

func NewDiscordBot(logger *slog.Logger, cfg *conf.Config, central centralproto.CentralServiceClient) (*DiscordBot, error) {
	if logger == nil {
		return nil, fmt.Errorf("nil logger not allowed")
	}

	if cfg == nil {
		return nil, fmt.Errorf("nil cfg not allowed")
	}

	client, err := disgo.New(
		cfg.BotToken,
		disgobot.WithGatewayConfigOpts(
			// set enabled intents
			gateway.WithIntents(
				gateway.IntentGuilds,
				gateway.IntentGuildMessages,
				gateway.IntentDirectMessages,
				gateway.IntentGuildInvites,
				gateway.IntentMessageContent,
			),
		),
		disgobot.WithLogger(logger),
	)
	if err != nil {
		return nil, err
	}

	bot := &DiscordBot{logger: logger, disc: client, central: central}

	return bot, nil
}
