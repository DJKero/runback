//! `SeaORM` Entity. Generated by sea-orm-codegen 0.6.0

use sea_orm::entity::prelude::*;
use serde::{Deserialize, Serialize};
use twilight_model::id::marker::UserMarker;

use crate::entity::IdWrapper;

#[derive(Clone, Debug, PartialEq, DeriveEntityModel, Serialize, Deserialize)]
#[sea_orm(table_name = "users")]
pub struct Model {
    #[sea_orm(primary_key, auto_increment = false)]
    pub user_id: Uuid,
    #[sea_orm(unique)]
    pub discord_user: Option<IdWrapper<UserMarker>>,
}

#[derive(Copy, Clone, Debug, EnumIter, DeriveRelation)]
pub enum Relation {
    #[sea_orm(has_many = "super::matchmaking_lobbies::Entity")]
    MatchmakingLobbies,
    #[sea_orm(has_many = "super::matchmaking_player_lobby::Entity")]
    MatchmakingPlayerLobby,
    #[sea_orm(has_many = "super::state::Entity")]
    State,
    #[sea_orm(has_many = "super::matchmaking_invitation::Entity")]
    MatchmakingInvitation,
    #[sea_orm(has_many = "super::matchmaking_player_invitation::Entity")]
    MatchmakingPlayerInvitation,
}

impl Related<super::matchmaking_lobbies::Entity> for Entity {
    fn to() -> RelationDef {
        Relation::MatchmakingLobbies.def()
    }
}

impl Related<super::matchmaking_player_lobby::Entity> for Entity {
    fn to() -> RelationDef {
        Relation::MatchmakingPlayerLobby.def()
    }
}

impl Related<super::state::Entity> for Entity {
    fn to() -> RelationDef {
        Relation::State.def()
    }
}

impl Related<super::matchmaking_invitation::Entity> for Entity {
    fn to() -> RelationDef {
        Relation::MatchmakingInvitation.def()
    }
}

impl Related<super::matchmaking_player_invitation::Entity> for Entity {
    fn to() -> RelationDef {
        Relation::MatchmakingPlayerInvitation.def()
    }
}

impl ActiveModelBehavior for ActiveModel {}
