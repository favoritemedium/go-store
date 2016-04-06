--
--  Tables needed by user.go and authuser.go.
--

CREATE TABLE `users` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `authby` int(10) unsigned NOT NULL,
  `email` varchar(255) NOT NULL,
  `pwhash` varchar(64) NOT NULL,
  `fullname` varchar(50) NOT NULL,
  `nametouse` varchar(50) NOT NULL,
  `isactive` tinyint(1) unsigned NOT NULL,
  `roles` int(10) unsigned NOT NULL,
  `created_at` int(10) unsigned NOT NULL,
  `updated_at` int(10) unsigned NOT NULL,
  `active_at` int(10) unsigned NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `users_signins` (
  `id` int(10) unsigned NOT NULL,
  `user_id` int(10) unsigned NOT NULL,
  `signedinat` int(10) unsigned NOT NULL,
  `ip` varchar(50) NOT NULL,
  `agent` varchar(50) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`),
  CONSTRAINT `users_signins_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `users_tokens` (
  `id` int(10) unsigned NOT NULL,
  `user_id` int(10) unsigned NOT NULL,
  `agent` varchar(50) NOT NULL,
  `auth` varchar(32) NOT NULL,
  `auth_expiry` int(10) unsigned NOT NULL,
  `refresh` varchar(32) NOT NULL,
  `refresh_expiry` int(10) unsigned NOT NULL,
  `active_at` int(10) unsigned NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `auth` (`auth`),
  KEY `user_id` (`user_id`) USING BTREE,
  KEY `refresh` (`refresh`),
  CONSTRAINT `users_tokens_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `users_verify` (
 `id` int(10) unsigned NOT NULL,
 `email` varchar(255) NOT NULL,
 `token` varchar(32) NOT NULL,
 `code` varchar(10) NOT NULL,
 `expiry` int(10) unsigned NOT NULL,
 `foruser_id` int(10) unsigned NOT NULL,
 PRIMARY KEY (`id`),
 UNIQUE KEY `token` (`token`),
 KEY `code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
