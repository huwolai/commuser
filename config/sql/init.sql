-- +migrate Up

-- 用户表
CREATE TABLE IF NOT EXISTS user(
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  open_id VARCHAR(255) COMMENT 'open_id 此ID来自于用户中心',
  app_id VARCHAR(255)  COMMENT '应用ID',
  nickname VARCHAR(255) COMMENT '昵称',
  username VARCHAR(255) COMMENT '用户名',
  email VARCHAR(255) COMMENT '邮箱',
  mobile VARCHAR(255) COMMENT '手机号',
  password VARCHAR(255) COMMENT '密码',
  status int COMMENT '用户状态 1.可用 0.不可用',
  is_payapi int COMMENT '是否开通支付API 0.否 1.是',
  create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  update_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间戳'
) CHARACTER SET utf8;
