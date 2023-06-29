package livekit_helper

import (
	"context"
	"fmt"
	"github.com/chenxyzl/glin/slog"
	"github.com/livekit/protocol/auth"
	"github.com/livekit/protocol/livekit"
	lksdk "github.com/livekit/server-sdk-go"
	"golang.org/x/exp/slices"
	"laiya/config"
	"sync"
	"time"
)

var lockIns = sync.Mutex{}
var roomClient *lksdk.RoomServiceClient

func GetClient() *lksdk.RoomServiceClient {
	if roomClient == nil {
		lockIns.Lock()
		defer lockIns.Unlock()
		if roomClient != nil {
			return roomClient
		}
		roomClient = lksdk.NewRoomServiceClient(config.Get().AConfig.Livekit.Host, config.Get().AConfig.Livekit.ApiKey, config.Get().AConfig.Livekit.ApiSecret)
	}
	return roomClient
}

func GetToken(uid uint64, gid uint64) (string, error) {
	at := GetClient().CreateToken()
	grant := &auth.VideoGrant{
		RoomJoin: true,
		Room:     fmt.Sprintf("%d", gid),
	}
	at.AddGrant(grant).
		SetIdentity(fmt.Sprintf("%d", uid)).
		SetValidFor(time.Hour)
	return at.ToJWT()
}

func GetVoicePlayers(gid uint64) (*livekit.ListParticipantsResponse, error) {
	ctx, _ := context.WithTimeout(context.Background(), 500*time.Millisecond)
	return GetClient().ListParticipants(ctx, &livekit.ListParticipantsRequest{
		Room: fmt.Sprintf("%d", gid),
	})
}

func RemoveVoicePlayer(uid uint64, gid uint64) (*livekit.RemoveParticipantResponse, error) {
	ctx, _ := context.WithTimeout(context.Background(), 500*time.Millisecond)
	return GetClient().RemoveParticipant(ctx, &livekit.RoomParticipantIdentity{
		Room:     fmt.Sprintf("%d", gid),
		Identity: fmt.Sprintf("%d", uid),
	})
}

func MuteGidVoicePlayerByGid(gid uint64, datas *livekit.ListParticipantsResponse, exclude []uint64, mute bool) {
	room := fmt.Sprintf("%d", gid)
	var excludeId []string
	for _, uid := range exclude {
		excludeId = append(excludeId, fmt.Sprintf("%d", uid))
	}
	for _, participant := range datas.Participants {
		//过滤掉
		if slices.Contains(excludeId, participant.GetIdentity()) {
			continue
		}
		for _, track := range participant.Tracks {
			ctx, _ := context.WithTimeout(context.Background(), 500*time.Millisecond)
			_, err := GetClient().MutePublishedTrack(ctx, &livekit.MuteRoomTrackRequest{
				Room:     room,
				Identity: participant.GetIdentity(),
				TrackSid: track.GetSid(),
				Muted:    mute,
			})
			if err != nil {
				slog.Errorf("live mute uid err ,gid:%v|uid:%v|tid:%v|err:%v", room, participant.GetIdentity(), track.GetSid(), err)
			}
		}
	}
	return
}

func MuteGidVoicePlayer(gid uint64, exclude []uint64, mute bool) (*livekit.ListParticipantsResponse, error) {
	rsp, err := GetVoicePlayers(gid)
	if err != nil {
		return nil, err
	}
	MuteGidVoicePlayerByGid(gid, rsp, exclude, mute)
	return rsp, nil
}

func MuteGidVoicePlayerByUid(gid uint64, datas *livekit.ListParticipantsResponse, uids []uint64, mute bool) {
	room := fmt.Sprintf("%d", gid)
	var targetList []string
	for _, uid := range uids {
		targetList = append(targetList, fmt.Sprintf("%d", uid))
	}
	for _, participant := range datas.Participants {
		//过滤掉
		if !slices.Contains(targetList, participant.GetIdentity()) {
			continue
		}
		for _, track := range participant.Tracks {
			ctx, _ := context.WithTimeout(context.Background(), 500*time.Millisecond)
			_, err := GetClient().MutePublishedTrack(ctx, &livekit.MuteRoomTrackRequest{
				Room:     room,
				Identity: participant.GetIdentity(),
				TrackSid: track.GetSid(),
				Muted:    mute,
			})
			if err != nil {
				slog.Errorf("live mute uid err ,gid:%v|uid:%v|tid:%v|err:%v", room, participant.GetIdentity(), track.GetSid(), err)
			}
		}
	}
	return
}
