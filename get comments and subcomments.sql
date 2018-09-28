SELECT C.id,C.content as c_content,C.nickname as c_nickname,C.creation_date as c_creation_date,
S.content as s_content,S.nickname as s_nickname,S.creation_date as s_creation_date
from comments C JOIN sub_comments S ON C.id=S.comment_id     
where C.article_id=1 and C.active =1