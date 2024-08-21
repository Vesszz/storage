package dbloader

import (
	"fmt"
	"go.uber.org/zap"
	"main/internal/models"
)

func (l *DBLoader) InsertCommentReaction(commentReaction models.ReactionComment) error {
	_, err := l.db.Exec("INSERT INTO reaction_comments(user_id, comment_id, type_id) VALUES ($1, $2, $3)",
		commentReaction.UserID, commentReaction.CommentID, commentReaction.TypeID)
	if err != nil {
		zap.S().Errorf("insert comment reaction: %v", err)
		return fmt.Errorf("insert comment reaction: %w", err)
	}
	return nil
}

func (l *DBLoader) DeleteCommentReaction(commentReaction models.ReactionComment) error {
	_, err := l.db.Exec("DELETE FROM reaction_files WHERE user_id = $1 AND file_id = $2",
		commentReaction.UserID, commentReaction.CommentID)
	if err != nil {
		zap.S().Errorf("delete comment reaction: %v", err)
		return fmt.Errorf("delete comment reaction: %w", err)
	}
	return nil
}
