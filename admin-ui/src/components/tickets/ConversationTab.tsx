'use client';

import { useState, useEffect } from 'react';
import { MessageSquare, Send, User, Clock } from 'lucide-react';
import { ticketsApi, type TicketComment } from '@/lib/api/tickets';

interface ConversationTabProps {
  ticketId: string;
}

export default function ConversationTab({ ticketId }: ConversationTabProps) {
  const [comments, setComments] = useState<TicketComment[]>([]);
  const [loading, setLoading] = useState(true);
  const [newComment, setNewComment] = useState('');

  useEffect(() => {
    loadComments();
  }, [ticketId]);

  const loadComments = async () => {
    try {
      setLoading(true);
      const response = await ticketsApi.getComments(ticketId);
      setComments(response.comments || []);
    } catch (error) {
      console.error('Failed to load comments:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleSendComment = async () => {
    if (!newComment.trim()) return;

    try {
      await ticketsApi.addComment(ticketId, {
        comment: newComment,
        comment_type: 'engineer',
      });
      setNewComment('');
      loadComments();
    } catch (error) {
      console.error('Failed to send comment:', error);
    }
  };

  if (loading) {
    return (
      <div className="bg-white rounded-lg shadow p-6">
        <div className="flex items-center justify-center py-12">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
        </div>
      </div>
    );
  }

  return (
    <div className="bg-white rounded-lg shadow">
      {/* Comments List */}
      <div className="p-6 space-y-4 max-h-[600px] overflow-y-auto">
        {comments.length === 0 ? (
          <div className="text-center py-12">
            <MessageSquare className="h-12 w-12 text-gray-300 mx-auto mb-4" />
            <p className="text-gray-500">No comments yet. Start the conversation!</p>
          </div>
        ) : (
          comments.map((comment) => (
            <div key={comment.id} className="flex space-x-3">
              <div className="flex-shrink-0">
                <div className="h-8 w-8 bg-blue-100 rounded-full flex items-center justify-center">
                  <User className="h-4 w-4 text-blue-600" />
                </div>
              </div>
              <div className="flex-1">
                <div className="flex items-center space-x-2 mb-1">
                  <span className="text-sm font-semibold text-gray-900">
                    {comment.author_name}
                  </span>
                  <span className="text-xs text-gray-500">
                    {new Date(comment.created_at).toLocaleString()}
                  </span>
                  <CommentTypeBadge type={comment.comment_type} />
                </div>
                <p className="text-sm text-gray-700">{comment.comment}</p>
              </div>
            </div>
          ))
        )}
      </div>

      {/* New Comment Input */}
      <div className="border-t border-gray-200 p-4">
        <div className="flex space-x-3">
          <textarea
            value={newComment}
            onChange={(e) => setNewComment(e.target.value)}
            placeholder="Add a comment..."
            className="flex-1 px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent resize-none"
            rows={3}
          />
          <button
            onClick={handleSendComment}
            disabled={!newComment.trim()}
            className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition disabled:opacity-50 disabled:cursor-not-allowed self-end"
          >
            <Send className="h-4 w-4" />
          </button>
        </div>
      </div>
    </div>
  );
}

function CommentTypeBadge({ type }: { type?: string }) {
  const badges = {
    customer: 'bg-green-100 text-green-700',
    engineer: 'bg-blue-100 text-blue-700',
    internal: 'bg-yellow-100 text-yellow-700',
    system: 'bg-gray-100 text-gray-700',
  };

  return (
    <span
      className={`inline-flex px-2 py-0.5 text-xs font-medium rounded ${
        badges[type as keyof typeof badges] || 'bg-gray-100 text-gray-700'
      }`}
    >
      {type}
    </span>
  );
}
