-- Create users table first
CREATE TABLE `users` (
  `id` INT(11) NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(255) NOT NULL,
  `kakao_id` VARCHAR(255) UNIQUE, -- kakao_id를 UNIQUE로 설정
  `skku_mail` VARCHAR(255),
  `phone` VARCHAR(255) NOT NULL,
  `usertype` ENUM('external', 'skkuin') NOT NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
);

-- Create skkuin table
CREATE TABLE `skkuin` (
  `id` INT(11) NOT NULL AUTO_INCREMENT,
  `skkuin_type` ENUM('student', 'professor', 'staff') NOT NULL,
  `department` VARCHAR(255) NOT NULL,
  `student_id` VARCHAR(255),
  `user_id` INT(11),
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE CASCADE
);

-- Create sessions table
CREATE TABLE `sessions` (
  `session_id` VARCHAR(255) NOT NULL, -- 고유 세션 ID
  `kakao_id` VARCHAR(255) NOT NULL,
  `original_url` TEXT NOT NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `expires_at` TIMESTAMP NOT NULL DEFAULT (CURRENT_TIMESTAMP + INTERVAL 10 MINUTE),
  PRIMARY KEY (`session_id`), 
  FOREIGN KEY (`kakao_id`) REFERENCES `users`(`kakao_id`) ON DELETE CASCADE
);

-- 만료된 세션 자동 삭제 이벤트 생성 (MySQL EVENT 스케줄러 사용)
DROP EVENT IF EXISTS delete_expired_sessions;
CREATE EVENT delete_expired_sessions
ON SCHEDULE EVERY 5 MINUTE
DO
  DELETE FROM sessions WHERE expires_at < NOW();


-- SET GLOBAL event_scheduler = ON; -- 이벤트 스케줄러 활성화 (mysql 1회성 시행)

