-- +migrate Up
alter table `user` add column flag varchar(50) DEFAULT '' COMMENT '标记';
alter table `user` add column json varchar(1000) DEFAULT '' COMMENT '附加数据';