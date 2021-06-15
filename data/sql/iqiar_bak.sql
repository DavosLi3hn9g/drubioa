/*
 Navicat Premium Data Transfer

 Source Server         : iqiar
 Source Server Type    : SQLite
 Source Server Version : 3030001
 Source Schema         : main

 Target Server Type    : SQLite
 Target Server Version : 3030001
 File Encoding         : 65001

*/

PRAGMA foreign_keys = false;

-- ----------------------------
-- Table structure for iq_auth
-- ----------------------------
DROP TABLE IF EXISTS "iq_auth";
CREATE TABLE "iq_auth"
(
	auth_id integer
		primary key autoincrement,
	ai_name varchar(255),
	password varchar(255),
	activated integer(1) default 0
);


-- ----------------------------
-- Table structure for iq_call_name
-- ----------------------------
DROP TABLE IF EXISTS "iq_call_name";
CREATE TABLE "iq_call_name" ("id" integer primary key autoincrement,"uid" integer,"call" varchar(255) );

-- ----------------------------
-- Records of iq_call_name
-- ----------------------------
BEGIN;
INSERT INTO "iq_call_name" VALUES (10, 1, '先生');
COMMIT;

-- ----------------------------
-- Table structure for iq_intentions
-- ----------------------------
DROP TABLE IF EXISTS "iq_intentions";
CREATE TABLE "iq_intentions" ("sid" integer primary key autoincrement,"title" varchar(255),"end" varchar(255),"level" integer,"hello" bool , "hits" integer);

-- ----------------------------
-- Records of iq_intentions
-- ----------------------------
BEGIN;
INSERT INTO "iq_intentions" VALUES (1, '闲聊', '', 100, 1, NULL);
INSERT INTO "iq_intentions" VALUES (2, '房屋中介', 1, 90, 0, NULL);
INSERT INTO "iq_intentions" VALUES (3, '银行贷款', 1, 80, 0, NULL);
INSERT INTO "iq_intentions" VALUES (4, '快递', 2, 100, 0, NULL);
INSERT INTO "iq_intentions" VALUES (5, '炒股', 1, 70, 0, NULL);
COMMIT;

-- ----------------------------
-- Table structure for iq_logs_call
-- ----------------------------
DROP TABLE IF EXISTS "iq_logs_call";
CREATE TABLE "iq_logs_call" ("id" integer primary key autoincrement,"content" varchar(255),"recording" varchar(255),"tel_from" varchar(255),"tel_to" varchar(255),"intention" varchar(255),"policy" varchar(255),"time_start" integer,"time_end" integer,"minute" integer , "text" varchar(255));

-- ----------------------------
-- Records of iq_logs_call
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for iq_logs_sms
-- ----------------------------
DROP TABLE IF EXISTS "iq_logs_sms";
CREATE TABLE "iq_logs_sms" ("id" integer primary key autoincrement,"text" varchar(255),"tel_from" varchar(255),"tel_to" varchar(255),"dateline" integer );

-- ----------------------------
-- Records of iq_logs_sms
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for iq_policies
-- ----------------------------
DROP TABLE IF EXISTS "iq_policies";
CREATE TABLE "iq_policies" ("id" integer primary key autoincrement,"title" varchar(255),"checked" varchar(255),"silent" bool , "hits" integer);

-- ----------------------------
-- Records of iq_policies
-- ----------------------------
BEGIN;
INSERT INTO "iq_policies" VALUES (1, '不需要，请不要再打这个电话！', 'sms|blacklist', 0, 3);
INSERT INTO "iq_policies" VALUES (2, '抱歉，其实我是机器人，我正在通知我的主人', 'push|sms|call', 0, 0);
INSERT INTO "iq_policies" VALUES (3, '好的，其实我是机器人，如果有需要主人会联系你', 'sms', 0, 0);
COMMIT;

-- ----------------------------
-- Table structure for iq_queries
-- ----------------------------
DROP TABLE IF EXISTS "iq_queries";
CREATE TABLE "iq_queries" ("id" integer primary key autoincrement,"sid" integer,"scores" integer,"query" varchar(255),"answer" varchar(255) , "mode" integer);

-- ----------------------------
-- Records of iq_queries
-- ----------------------------
BEGIN;
INSERT INTO "iq_queries" VALUES (1, 1, 0, '喂,喂。', '在听,请说,什么事？', 1);
INSERT INTO "iq_queries" VALUES (2, 1, 0, '你好,您好', '有什么事？', 1);
INSERT INTO "iq_queries" VALUES (3, 2, 50, '出租,出售,房产,商户,买房', '是房屋出售还是出租？', 1);
INSERT INTO "iq_queries" VALUES (4, 3, 50, '贷款,资金,信贷', '你是说贷款吗？不需要！', 1);
INSERT INTO "iq_queries" VALUES (5, 3, 30, '利息', '利息太高，不考虑！', 1);
INSERT INTO "iq_queries" VALUES (6, 4, 50, '楼下,下楼', '', 1);
INSERT INTO "iq_queries" VALUES (7, 4, 50, '快递', '您是到楼下了吗？', 1);
INSERT INTO "iq_queries" VALUES (8, 4, 50, '包裹', '', 1);
INSERT INTO "iq_queries" VALUES (9, 4, 50, '取餐', '', 1);
INSERT INTO "iq_queries" VALUES (10, 4, 50, '取一下', '', 1);
INSERT INTO "iq_queries" VALUES (11, 1, 30, '需要吗,需求,是否需要,需不需要,有兴趣吗', '', 1);
INSERT INTO "iq_queries" VALUES (12, 1, 1, '接电话,我找,找一下,联系他', '有什么事？', 1);
INSERT INTO "iq_queries" VALUES (13, 1, 0, '请问您是,请问是,请问你是', '有什么事？', 1);
INSERT INTO "iq_queries" VALUES (14, 1, 30, '我们是做,这边是做,我们公司', '', 1);
INSERT INTO "iq_queries" VALUES (15, 1, 5, '听到吗,在听吗,在不在,听不到,没听清,再说一遍,没有太听清', '请说', 1);
INSERT INTO "iq_queries" VALUES (16, 1, 0, '机器人吗,机器人吧', '你，猜', 1);
INSERT INTO "iq_queries" VALUES (17, 1, 50, '专业办理', '', 1);
INSERT INTO "iq_queries" VALUES (18, 3, 30, '优质客户', '您是从哪里得到的我的电话号码的？', 1);
INSERT INTO "iq_queries" VALUES (19, 3, 50, '随借随还,当天到账', '不需要贷款', 1);
INSERT INTO "iq_queries" VALUES (20, 3, 30, '授信,额度,银行,信用社', '', 1);
INSERT INTO "iq_queries" VALUES (21, 2, 1, '投资吗', '投资什么？', 1);
INSERT INTO "iq_queries" VALUES (22, 5, 30, '散户,炒股', '我不炒股', 1);
INSERT INTO "iq_queries" VALUES (23, 5, 30, '股票交流', '我不炒股', 1);
INSERT INTO "iq_queries" VALUES (24, 1, 30, '客户经理,客服经理', '您是哪个公司的？', 1);
INSERT INTO "iq_queries" VALUES (25, 1, 0, '答谢老用户,出席会议', '您是从哪里得到的我的电话号码的？', 1);
INSERT INTO "iq_queries" VALUES (26, 4, 50, '开门吧,已经到了', '请稍等，我正在通知我的主人。', 1);
INSERT INTO "iq_queries" VALUES (27, 2, 30, '售楼部', '抱歉，没钱买房', 1);
COMMIT;

-- ----------------------------
-- Table structure for iq_settings
-- ----------------------------
DROP TABLE IF EXISTS "iq_settings";
CREATE TABLE "iq_settings" ("key" varchar(255),"value" varchar(255) , "history" varchar(255), PRIMARY KEY ("key"));

-- ----------------------------
-- Records of iq_settings
-- ----------------------------
BEGIN;
INSERT INTO "iq_settings" VALUES ('talk_prologue', '你好，有什么事？', '您好，找谁？|您好，请说话。|有什么事？');
INSERT INTO "iq_settings" VALUES ('tts_voice', 'Xiaowei', '');
INSERT INTO "iq_settings" VALUES ('tts_speech', 50, '');
INSERT INTO "iq_settings" VALUES ('tts_pitch', 0, '');
INSERT INTO "iq_settings" VALUES ('tts_volume', 50, '');
INSERT INTO "iq_settings" VALUES ('tty_call', '/dev/ttyAMA0', '');
INSERT INTO "iq_settings" VALUES ('tty_usb', '/dev/ttyUSB4', '');
INSERT INTO "iq_settings" VALUES ('aliyun_ak_id', '', '');
INSERT INTO "iq_settings" VALUES ('aliyun_ak_secret', '', '');
INSERT INTO "iq_settings" VALUES ('aliyun_isi_appkey', '', '');
INSERT INTO "iq_settings" VALUES ('aliyun_bucket_name', '', '');
INSERT INTO "iq_settings" VALUES ('upload_oss', 'true', '');
INSERT INTO "iq_settings" VALUES ('unit_price', 0.2, '');
COMMIT;

-- ----------------------------
-- Table structure for iq_users
-- ----------------------------
DROP TABLE IF EXISTS "iq_users";
CREATE TABLE "iq_users" ("uid" integer primary key autoincrement,"user_name" varchar(255),"tel" varchar(255),"email" varchar(255),"push" varchar(255),"lev" integer , "hits" integer);

-- ----------------------------
-- Records of iq_users
-- ----------------------------
BEGIN;
INSERT INTO "iq_users" VALUES (1, '我自己', '', '', '', 0, NULL);
INSERT INTO "iq_users" VALUES (2, '其他人', '', '', '', 0, NULL);
COMMIT;

-- ----------------------------
-- Auto increment value for iq_auth
-- ----------------------------
UPDATE "main"."sqlite_sequence" SET seq = 1 WHERE name = 'iq_auth';

-- ----------------------------
-- Auto increment value for iq_call_name
-- ----------------------------
UPDATE "main"."sqlite_sequence" SET seq = 10 WHERE name = 'iq_call_name';

-- ----------------------------
-- Auto increment value for iq_intentions
-- ----------------------------
UPDATE "main"."sqlite_sequence" SET seq = 12 WHERE name = 'iq_intentions';

-- ----------------------------
-- Auto increment value for iq_logs_call
-- ----------------------------
UPDATE "main"."sqlite_sequence" SET seq = 106 WHERE name = 'iq_logs_call';

-- ----------------------------
-- Auto increment value for iq_logs_sms
-- ----------------------------
UPDATE "main"."sqlite_sequence" SET seq = 207 WHERE name = 'iq_logs_sms';

-- ----------------------------
-- Auto increment value for iq_policies
-- ----------------------------
UPDATE "main"."sqlite_sequence" SET seq = 3 WHERE name = 'iq_policies';

-- ----------------------------
-- Auto increment value for iq_queries
-- ----------------------------
UPDATE "main"."sqlite_sequence" SET seq = 79 WHERE name = 'iq_queries';

-- ----------------------------
-- Auto increment value for iq_users
-- ----------------------------
UPDATE "main"."sqlite_sequence" SET seq = 3 WHERE name = 'iq_users';

PRAGMA foreign_keys = true;
