CREATE TABLE `users` (
  `id` int(10) UNSIGNED NOT NULL,
  `allowedauth` int(10) UNSIGNED NOT NULL,
  `email` varchar(255) NOT NULL,
  `pwhash` varchar(64) NOT NULL,
  `fullname` varchar(50) NOT NULL,
  `nametouse` varchar(50) NOT NULL,
  `isactive` tinyint(1) NOT NULL,
  `isadmin` tinyint(1) NOT NULL,
  `issuperuser` tinyint(1) NOT NULL,
  `createdat` int(10) UNSIGNED NOT NULL,
  `updatedat` int(10) UNSIGNED NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

ALTER TABLE `users`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `email` (`email`);
