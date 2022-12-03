use crate::entity::prelude::*;
use crate::migration::deprecated::*;
use sea_orm_migration::prelude::*;

use super::deprecated::*;

#[derive(DeriveMigrationName)]
pub struct Migration;

#[async_trait::async_trait]
impl MigrationTrait for Migration {
    async fn up(&self, manager: &SchemaManager) -> Result<(), DbErr> {
        // let stmt = Statement::from_string(manager.get_database_backend(), "BEGIN".to_owned());
        // manager.get_connection().execute(stmt).await?;

        manager
            .alter_table(
                Table::alter()
                    .table(Game)
                    .add_column_if_not_exists(
                        ColumnDef::new(game::Column::SteamAppId)
                            .big_integer()
                            .not_null()
                            .unique_key(),
                    )
                    .rename_column(game::Column::Name, game::Column::GameName)
                    .to_owned(),
            )
            .await?;

        manager
            .drop_table(
                Table::drop()
                    .if_exists()
                    .table(GameCharacter)
                    .to_owned(),
            )
            .await?;
        
        manager
            .create_table(
                Table::create()
                    .table(Character)
                    .if_not_exists()
                    .col(
                        ColumnDef::new(character::Column::Id)
                            .big_integer()
                            .auto_increment()
                            .primary_key()
                            .not_null()
                            .unique_key(),
                    )
                    .col(
                        ColumnDef::new(character::Column::GameId)
                            .uuid()
                            .not_null(),
                    )
                    .foreign_key(
                        ForeignKey::create()
                            .name("FK_game_id")
                            .from(Character, character::Column::GameId)
                            .to(Game, game::Column::Id)
                            .on_delete(ForeignKeyAction::Cascade)
                            .on_update(ForeignKeyAction::Cascade)
                    )
                    .col(
                        ColumnDef::new(character::Column::CharacterName)
                            .string_len(80)
                            .not_null(),
                    )
                    .to_owned(),
            )
            .await?;

        manager
            .create_table(
                Table::create()
                    .table(Characters)
                    .if_not_exists()
                    .col(
                        ColumnDef::new(characters::Column::Id)
                            .uuid()
                            .primary_key()
                            .not_null()
                            .unique_key(),
                    )
                    .col(
                        ColumnDef::new(characters::Column::GameId)
                            .uuid()
                            .not_null(),
                    )
                    .foreign_key(
                        ForeignKey::create()
                            .name("FK_game_id")
                            .from(Characters, characters::Column::GameId)
                            .to(Game, game::Column::Id)
                            .on_delete(ForeignKeyAction::Cascade)
                            .on_update(ForeignKeyAction::Cascade)
                    )
                    .col(
                        ColumnDef::new(characters::Column::Characters)
                            .array("bigint".to_string())
                            .not_null(),
                    )
                    .to_owned(),
            )
            .await?;
        
        manager
            .create_table(
                Table::create()
                    .table(DuelMode)
                    .if_not_exists()
                    .col(
                        ColumnDef::new(duel_mode::Column::Id)
                            .uuid()
                            .primary_key()
                            .not_null()
                            .unique_key(),
                    )
                    .col(
                        ColumnDef::new(duel_mode::Column::GameId)
                            .uuid()
                            .not_null(),
                    )
                    .foreign_key(
                        ForeignKey::create()
                            .name("FK_game_id")
                            .from(DuelMode, duel_mode::Column::GameId)
                            .to(Game, game::Column::Id)
                            .on_delete(ForeignKeyAction::Cascade)
                            .on_update(ForeignKeyAction::Cascade)
                    )
                    .col(
                        ColumnDef::new(duel_mode::Column::ModeName)
                            .string_len(80)
                            .not_null(),
                    )
                    .col(
                        ColumnDef::new(duel_mode::Column::ModeRules)
                            .string_len(80)
                            .not_null(),
                    )
            )
            .await?;
        
        manager
            .create_table(
                Table::create()
                    .table(DuelRating)
                    .if_not_exists()
                    .col(
                        ColumnDef::new(duel_rating::Column::Id)
                            .uuid()
                            .primary_key()
                            .not_null()
                            .unique_key(),
                    )
                    .col(
                        ColumnDef::new(duel_rating::Column::UserId)
                            .uuid()
                            .not_null(),
                    )
                    .foreign_key(
                        ForeignKey::create()
                            .name("FK_user_id")
                            .from(DuelRating, duel_rating::Column::UserId)
                            .to(Users, users::Column::UserId)
                            .on_delete(ForeignKeyAction::Cascade)
                            .on_update(ForeignKeyAction::Cascade)
                    )
                    .col(
                        ColumnDef::new(duel_rating::Column::ModeId)
                            .uuid()
                            .not_null(),
                    )
                    .foreign_key(
                        ForeignKey::create()
                            .name("FK_duel_mode_id")
                            .from(DuelRating, duel_rating::Column::ModeId)
                            .to(DuelMode, duel_mode::Column::Id)
                            .on_delete(ForeignKeyAction::Cascade)
                            .on_update(ForeignKeyAction::Cascade)
                    )
                    .col(
                        ColumnDef::new(duel_rating::Column::CharactersId)
                            .uuid()
                            .not_null(),
                    )
                    .foreign_key(
                        ForeignKey::create()
                            .name("FK_characters_id")
                            .from(DuelRating, duel_rating::Column::CharactersId)
                            .to(Characters, characters::Column::Id)
                            .on_delete(ForeignKeyAction::Cascade)
                            .on_update(ForeignKeyAction::Cascade)
                    )
                    .col(
                        ColumnDef::new(duel_rating::Column::Rating)
                            .json()
                            .not_null(),
                    )
                    .col(
                        ColumnDef::new(duel_rating::Column::LastUpdate)
                            .timestamp_with_time_zone()
                            .not_null(),
                    )
                    .to_owned(),
            )
            .await?;
        
        manager
            .create_table(
                Table::create()
                    .table(DuelMatch)
                    .if_not_exists()
                    .col(
                        ColumnDef::new(duel_match::Column::Id)
                            .uuid()
                            .primary_key()
                            .not_null()
                            .unique_key(),
                    )
                    .col(
                        ColumnDef::new(duel_match::Column::UserId)
                            .uuid()
                            .not_null(),
                    )
                    .foreign_key(
                        ForeignKey::create()
                            .name("FK_user_id")
                            .from(DuelMatch, duel_match::Column::UserId)
                            .to(Users, users::Column::UserId)
                            .on_delete(ForeignKeyAction::Cascade)
                            .on_update(ForeignKeyAction::Cascade)
                    )
                    .col(
                        ColumnDef::new(duel_match::Column::OpponentId)
                            .uuid()
                            .not_null(),
                    )
                    .foreign_key(
                        ForeignKey::create()
                            .name("FK_opponent_id")
                            .from(DuelMatch, duel_match::Column::OpponentId)
                            .to(Users, users::Column::UserId)
                            .on_delete(ForeignKeyAction::Cascade)
                            .on_update(ForeignKeyAction::Cascade)
                    )
                    .col(
                        ColumnDef::new(duel_match::Column::ModeId)
                            .uuid()
                            .not_null(),
                    )
                    .foreign_key(
                        ForeignKey::create()
                            .name("FK_duel_mode_id")
                            .from(DuelMatch, duel_match::Column::ModeId)
                            .to(DuelMode, duel_mode::Column::Id)
                            .on_delete(ForeignKeyAction::Cascade)
                            .on_update(ForeignKeyAction::Cascade)
                    )
                    .col(
                        ColumnDef::new(duel_match::Column::CharactersId)
                            .uuid()
                            .not_null(),
                    )
                    .foreign_key(
                        ForeignKey::create()
                            .name("FK_characters_id")
                            .from(DuelMatch, duel_match::Column::CharactersId)
                            .to(Characters, characters::Column::Id)
                            .on_delete(ForeignKeyAction::Cascade)
                            .on_update(ForeignKeyAction::Cascade)
                    )
                    .col(
                        ColumnDef::new(duel_match::Column::Wins)
                            .integer()
                            .not_null(),
                    )
                    .col(
                        ColumnDef::new(duel_match::Column::Loses)
                            .integer()
                            .not_null(),
                    )
                    .col(
                        ColumnDef::new(duel_match::Column::OldRating)
                            .json()
                            .not_null(),
                    )
                    .col(
                        ColumnDef::new(duel_match::Column::NewRating)
                            .json()
                            .not_null(),
                    )
                    .col(
                        ColumnDef::new(duel_match::Column::CreationTimestamp)
                            .timestamp_with_time_zone()
                            .not_null(),
                    )
                    .to_owned(),
            )
            .await?;
        /*
        manager
            .create_table(
                Table::create()
                    .table(Team)
                    .if_not_exists()
                    .col(
                        ColumnDef::new(team::Column::Id)
                            .uuid()
                            .primary_key()
                            .not_null()
                            .unique_key(),
                    )
                    .col(
                        ColumnDef::new(team::Column::TeamName)
                            .string_len(80)
                            .not_null(),
                    )
                    .col(
                        ColumnDef::new(team::Column::CreationTimestamp)
                            .timestamp_with_time_zone()
                            .not_null(),
                    )
                    .to_owned(),
            )
            .await?;

        manager
            .create_table(
                Table::create()
                    .table(TeamRule)
                    .if_not_exists()
                    .col(
                        ColumnDef::new(team_rule::Column::Id)
                            .uuid()
                            .primary_key()
                            .not_null()
                            .unique_key(),
                    )
                    .col(
                        ColumnDef::new(team_rule::Column::MinTeams)
                            .integer()
                            .not_null(),
                    )
                    .col(
                        ColumnDef::new(team_rule::Column::MaxTeams)
                            .integer()
                            .not_null(),
                    )
                    .col(
                        ColumnDef::new(team_rule::Column::MinTeamSize)
                            .integer()
                            .not_null(),
                    )
                    .col(
                        ColumnDef::new(team_rule::Column::MaxTeamSize)
                            .integer()
                            .not_null(),
                    )
                    .to_owned(),
            )
            .await?;
        
        manager
            .create_table(
                Table::create()
                    .table(TeamMode)
                    .if_not_exists()
                    .col(
                        ColumnDef::new(team_mode::Column::Id)
                            .uuid()
                            .primary_key()
                            .not_null()
                            .unique_key(),
                    )
                    .col(
                        ColumnDef::new(team_mode::Column::GameId)
                            .uuid()
                            .not_null(),
                    )
                    .foreign_key(
                        ForeignKey::create()
                            .name("FK_game_id")
                            .from(TeamMode, team_mode::Column::GameId)
                            .to(Game, game::Column::Id)
                            .on_delete(ForeignKeyAction::Cascade)
                            .on_update(ForeignKeyAction::Cascade)
                    )
                    .col(
                        ColumnDef::new(team_mode::Column::ModeName)
                            .string_len(80)
                            .not_null(),
                    )
                    .col(
                        ColumnDef::new(team_mode::Column::ModeRules)
                            .string_len(80)
                            .not_null(),
                    )
                    .col(
                        ColumnDef::new(team_mode::Column::TeamRuleId)
                            .uuid(),
                    )
                    .foreign_key(
                        ForeignKey::create()
                            .name("FK_team_rule_id")
                            .from(DuelMode, team_mode::Column::TeamRuleId)
                            .to(TeamRule, team_rule::Column::Id)
                            .on_delete(ForeignKeyAction::Cascade)
                            .on_update(ForeignKeyAction::Cascade)
                    )
                    .to_owned(),
            )
            .await?;
        
        manager
            .create_table(
                Table::create()
                    .table(TeamRating)
                    .if_not_exists()
                    .col(
                        ColumnDef::new(team_rating::Column::Id)
                            .uuid()
                            .primary_key()
                            .not_null()
                            .unique_key(),
                    )
                    .col(
                        ColumnDef::new(team_rating::Column::TeamId)
                            .uuid()
                            .not_null(),
                    )
                    .foreign_key(
                        ForeignKey::create()
                            .name("FK_team_id")
                            .from(TeamRating, team_rating::Column::TeamId)
                            .to(Team, team::Column::Id)
                            .on_delete(ForeignKeyAction::Cascade)
                            .on_update(ForeignKeyAction::Cascade)
                    )
                    .col(
                        ColumnDef::new(team_rating::Column::ModeId)
                            .uuid()
                            .not_null(),
                    )
                    .foreign_key(
                        ForeignKey::create()
                            .name("FK_team_mode_id")
                            .from(TeamRating, team_rating::Column::ModeId)
                            .to(TeamMode, team_mode::Column::Id)
                            .on_delete(ForeignKeyAction::Cascade)
                            .on_update(ForeignKeyAction::Cascade)
                    )
                    .col(
                        ColumnDef::new(team_rating::Column::Rating)
                            .json()
                            .not_null(),
                    )
                    .col(
                        ColumnDef::new(team_rating::Column::LastUpdate)
                            .timestamp_with_time_zone()
                            .not_null(),
                    )
                    .to_owned(),
            )
            .await?;

        manager
            .create_table(
                Table::create()
                    .table(TeamMatch)
                    .if_not_exists()
                    .col(
                        ColumnDef::new(team_match::Column::Id)
                            .uuid()
                            .primary_key()
                            .not_null()
                            .unique_key(),
                    )
                    .col(
                        ColumnDef::new(team_match::Column::UserId)
                            .uuid()
                            .not_null(),
                    )
                    .foreign_key(
                        ForeignKey::create()
                            .name("FK_user_id")
                            .from(TeamMatch, team_match::Column::UserId)
                            .to(Users, users::Column::UserId)
                            .on_delete(ForeignKeyAction::Cascade)
                            .on_update(ForeignKeyAction::Cascade)
                    )
                    .col(
                        ColumnDef::new(team_match::Column::OpponentId)
                            .uuid()
                            .not_null(),
                    )
                    .foreign_key(
                        ForeignKey::create()
                            .name("FK_opponent_id")
                            .from(TeamMatch, team_match::Column::OpponentId)
                            .to(Users, users::Column::UserId)
                            .on_delete(ForeignKeyAction::Cascade)
                            .on_update(ForeignKeyAction::Cascade)
                    )
                    .col(
                        ColumnDef::new(team_match::Column::ModeId)
                            .uuid()
                            .not_null(),
                    )
                    .foreign_key(
                        ForeignKey::create()
                            .name("FK_team_mode_id")
                            .from(TeamMatch, team_match::Column::ModeId)
                            .to(TeamMode, team_mode::Column::Id)
                            .on_delete(ForeignKeyAction::Cascade)
                            .on_update(ForeignKeyAction::Cascade)
                    )
                    .col(
                        ColumnDef::new(team_match::Column::Wins)
                            .integer()
                            .not_null(),
                    )
                    .col(
                        ColumnDef::new(team_match::Column::Loses)
                            .integer()
                            .not_null(),
                    )
                    .col(
                        ColumnDef::new(team_match::Column::OldRating)
                            .json()
                            .not_null(),
                    )
                    .col(
                        ColumnDef::new(team_match::Column::NewRating)
                            .json()
                            .not_null(),
                    )
                    .col(
                        ColumnDef::new(team_match::Column::CreationTimestamp)
                            .timestamp_with_time_zone()
                            .not_null(),
                    )
                    .to_owned(),
            )
            .await?;*/

        // let stmt = Statement::from_string(manager.get_database_backend(), "COMMIT".to_owned());
        // manager.get_connection().execute(stmt).await?;

        Ok(())
    }

    async fn down(&self, manager: &SchemaManager) -> Result<(), DbErr> {
        // Replace the sample below with your own migration scripts



        Ok(())
    }
}
