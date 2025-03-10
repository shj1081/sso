-- 삭제할 테이블의 순서를 고려하여 외래 키 제약 조건 문제 방지
DROP EVENT IF EXISTS delete_expired_sessions;
DROP TABLE IF EXISTS `sessions`;
DROP TABLE IF EXISTS `skkuin`;
DROP TABLE IF EXISTS `users`;