package common

// 登录用户的Session
var LoginSessionName = "LoginSessionName"

const FindRelationDetailedByWhereSql = `SELECT
		ps.id as relation_id,
		ps.user_id,
		ps.project_id,
		s.id as server_id,
		s.name AS server_name,
		s.ip_address,
		s.port,
		s.user_name,
		s.private_key,
		s.tag as server_tag,
		s.status,
		p.repo_name,
		p.branch_name,
		p.tag as project_tag,
		p.shell,
		p.status as project_status,
		p.key,
		p.secure,
		p.repo_type
	FROM
		eo_ps_relation AS ps
		LEFT JOIN eo_server AS s ON ps.server_id = s.id
		LEFT JOIN eo_project AS p ON ps.project_id = p.id
	WHERE
		1 = 1 and s.status = 0 `

const QueryByProjetIdSql = `SELECT
		s.*, ps.user_id,
		ps.create_time,
		ps.project_id,
		ps.id
	FROM
		eo_ps_relation as ps
		LEFT JOIN eo_server as s ON ps.server_id = s .id
	WHERE
		ps.project_id = ?
		AND  s.create_at = ?
	ORDER BY
		ps.id DESC`

const SearchServerByKeyword = `SELECT 
		* 
	FROM 
		eo_server as s 
	WHERE 
		s.create_at = ? AND (s.name LIKE ? OR s.tag LIKE ?)`
