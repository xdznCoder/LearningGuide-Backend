package logic

import (
	"LearningGuide/post_srv/internal/model"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	proto "LearningGuide/post_srv/.PostProto"
	"LearningGuide/post_srv/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetNoticeListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetNoticeListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetNoticeListLogic {
	return &GetNoticeListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetNoticeListLogic) GetNoticeList(in *proto.NoticeFilterRequest) (*proto.NoticeListResponse, error) {
	var count int64
	var notices []model.Notice

	filter := l.svcCtx.DB.Model(&model.Notice{}).Where("owner_id = ? AND owner_id != user_id", in.UserId)

	if in.Type != 0 {
		filter = filter.Where("type = ?", in.Type)
	}

	filter.Count(&count)
	result := filter.Scopes(model.Paginate(int(in.PageNum), int(in.PageSize))).Find(&notices)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	var respList []*proto.NoticeInfoResponse

	var ids []int32

	for _, n := range notices {
		ids = append(ids, n.ID)
		respList = append(respList, &proto.NoticeInfoResponse{
			Id:      n.ID,
			UserId:  n.UserId,
			OwnerId: n.OwnerId,
			Type:    n.Type,
			PostId:  n.PostId,
		})
	}

	err := l.svcCtx.DB.Model(&model.Notice{}).Where("Id in (?)", ids).Update("is_read", true).Error

	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &proto.NoticeListResponse{
		Total: count,
		Data:  respList,
	}, nil
}
