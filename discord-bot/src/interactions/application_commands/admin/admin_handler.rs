use std::sync::Arc;

use twilight_model::application::command::{
    BaseCommandOptionData, ChannelCommandOptionData, CommandOption, CommandType,
};
use twilight_model::channel::ChannelType;
use twilight_model::guild::Permissions;
use twilight_util::builder::command::{CommandBuilder, SubCommandBuilder, SubCommandGroupBuilder};

use crate::interactions::application_commands::{
    ApplicationCommandData, CommandGroupDescriptor, InteractionHandler, MessageComponentData,
};

use crate::interactions::application_commands::CommonUtilities;

use super::mm_settings_handler::MatchmakingSettingsHandler;

pub struct AdminCommandHandler {
    // utils: Arc<ApplicationCommandUtilities>,
    matchmaking_settings_handler: MatchmakingSettingsHandler,
    // matchmaking_panels_handler: MatchmakingPanelsHandler,
}

#[async_trait]
impl InteractionHandler for AdminCommandHandler {
    fn describe(&self) -> CommandGroupDescriptor {
        let builder = CommandBuilder::new(
            "admin",
            "Admin configuration and management settings",
            CommandType::ChatInput,
        )
        .dm_permission(false)
        .default_member_permissions(Permissions::MANAGE_GUILD)
        // .option(SubCommandBuilder::new(
        //     "matchmaking-panels".into(),
        //     "Add, edit, and remove matchmaking panels in your guild".into(),
        // ))
        .option(
            SubCommandGroupBuilder::new(
                "matchmaking-settings",
                "Shows the matchmaking settings panel",
            )
            .subcommands([
                SubCommandBuilder::new("admin-role", "Set which users can act as admins").option(
                    CommandOption::Role(BaseCommandOptionData {
                        name: "role".to_string(),
                        description: "The admin role (to disable, set to empty)".to_string(),
                        description_localizations: None,
                        name_localizations: None,
                        required: false,
                    }),
                ),
                SubCommandBuilder::new(
                    "matchmaking-channel",
                    "Set the default matchmaking channel",
                )
                .option(CommandOption::Channel(ChannelCommandOptionData {
                    name: "channel".to_string(),
                    description: "The matchmaking channel (to disable, set to empty)".to_string(),
                    channel_types: vec![ChannelType::GuildText],
                    description_localizations: None,
                    name_localizations: None,
                    required: false,
                })),
            ]),
        );

        let command = builder.build();
        CommandGroupDescriptor {
            name: "admin",
            description: "Tools for admins",
            commands: Box::new([command]),
        }
    }

    async fn process_command(&self, data: Box<ApplicationCommandData>) -> anyhow::Result<()> {
        let options = &data.command.options;

        if options.len() != 1 {
            return Err(anyhow!("Expected extra options when calling the top-level admin command handler. Number of arguments found: {}", options.len()));
        }

        let option = &options[0];

        let sub_command_name = options
            .get(0)
            .ok_or_else(|| anyhow!("Could not get first admin subcommand"))?
            .name
            .as_str();
        match sub_command_name {
            "matchmaking-settings" => {
                self.matchmaking_settings_handler
                    .process_command(data)
                    .await?;
            }
            "matchmaking-panels" => {
                // self.matchmaking_panels_handler
                //     .process_command(data)
                //     .await?;
                return Err(anyhow!("Panels are no longer supported"));
            }
            _ => {
                debug!(name = %option.name.as_str(), "Unknown admin subcommand option");
                return Err(anyhow!("Unknown admin subcommand option"));
            }
        }

        Ok(())
    }

    async fn process_autocomplete(&self, _data: Box<ApplicationCommandData>) -> anyhow::Result<()> {
        unreachable!()
    }

    async fn process_modal(&self, _data: Box<ApplicationCommandData>) -> anyhow::Result<()> {
        todo!("Admin handler does not currently process modals")
    }

    async fn process_component(&self, mut data: Box<MessageComponentData>) -> anyhow::Result<()> {
        debug!(leftovers = ?data.action, "Custom ID leftovers");

        if let Some((action, field)) = data.action.split_once(':') {
            match action {
                "settings" => {
                    data.action = field.to_string();
                    self.matchmaking_settings_handler
                        .process_component(data)
                        .await?;
                    return Ok(());
                }
                "panels" => {
                    // data.action = field.to_string();
                    // self.matchmaking_panels_handler
                    //     .process_component(data)
                    //     .await?;
                    // return Ok(());
                    return Err(anyhow!("Panels are no longer supported"));
                }
                _ => {
                    return Err(anyhow!("Unhandled action: {}", data.action));
                }
            }
        } else {
            return Err(anyhow!("Action did not match the format \"action:field\""));
        }
    }
}

impl AdminCommandHandler {
    pub fn new(utils: Arc<CommonUtilities>) -> Self {
        Self {
            matchmaking_settings_handler: MatchmakingSettingsHandler::new(utils),
            // matchmaking_panels_handler: MatchmakingPanelsHandler::new(utils.clone()),
            // utils,
        }
    }
}
