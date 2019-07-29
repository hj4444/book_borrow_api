SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for book
-- ----------------------------
DROP TABLE IF EXISTS `book`;
CREATE TABLE `book`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `description` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `url` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `ctime` timestamp(0) NULL DEFAULT CURRENT_TIMESTAMP,
  `mtime` timestamp(0) NULL DEFAULT NULL,
  `status` tinyint(2) NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 4 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of book
-- ----------------------------
INSERT INTO `book` VALUES ('tensorflow', '深度学习框架', 'http://img3m1.ddimg.cn/0/27/25224111-1_h_6.jpg.webp', '2019-01-16 10:58:27', '2019-01-16 10:55:53', 0);
INSERT INTO `book` VALUES ('蜜汁炖鱿鱼', '《亲爱的热爱的》原著小说，随书送杨紫、李现明信片，新增“婚礼+蜜月”高甜番外完整版，剧方唯.1正版授权', 'http://img3m6.ddimg.cn/9/10/26923356-1_b_1.jpg', '2019-07-25 09:38:31', NULL, 0);
INSERT INTO `book` VALUES ('只能陪你走一程', '蕊希2019新作，当当100%独家亲笔签名+随机彩蛋特签，预售期间，更有蕊希暖心福利大奖等你拿', 'http://img3m4.ddimg.cn/60/23/27904794-1_b_5.jpg', '2019-07-25 09:40:25', NULL, 0);

SET FOREIGN_KEY_CHECKS = 1;
