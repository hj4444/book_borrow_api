SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for user_book_map
-- ----------------------------
DROP TABLE IF EXISTS `user_book_map`;
CREATE TABLE `user_book_map`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL,
  `book_id` int(11) NOT NULL,
  `borrow_begin_date` timestamp(0) NOT NULL ON UPDATE CURRENT_TIMESTAMP(0),
  `borrow_end_date` timestamp(0) NOT NULL ON UPDATE CURRENT_TIMESTAMP(0),
  `ctime` timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `mtime` timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 30 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

SET FOREIGN_KEY_CHECKS = 1;
