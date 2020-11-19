<?php

return [
    'default' => 'things-db',
    'migrations' => 'migrations',
    'connections' => [
        'mysql' => [
            'driver'    => 'mysql',
            'host'      => env('DB_HOST'),
            'database'  => env('DB_DATABASE'),
            'port'      => env('DB_PORT'),
            'username'  => env('DB_USERNAME'),
            'password'  => env('DB_PASSWORD'),
            'charset'   => 'utf8',
            'collation' => 'utf8_unicode_ci',
            'prefix'    => '',
            'strict'    => false,
        ],
        'mysql_postal_code' => [
            'driver'    => 'mysql',
            'host'      => env('DB_HOST', 'localhost'),
            'database'  => env('DB_DATABASE_POSTAL_CODE'),
            'port'      => env('DB_PORT'),
            'username'  => env('DB_USERNAME'),
            'password'  => env('DB_PASSWORD'),
            'charset'   => 'utf8',
            'collation' => 'utf8_unicode_ci',
            'prefix'    => '',
            'strict'    => false,
        ],
        'users-db' => [
            'driver'   => 'pgsql',
            'host'     => env('DB_USERS_HOST', 'database_p'),
            'database' => env('DB_USERS_DATABASE', 'dockerApp'), // This seems to be ignored
            'port'     => env('DB_USERS_PORT', 5432),
            'username' => env('DB_USERS_USERNAME', 'postgres'),
            'password' => env('DB_USERS_PASSWORD', 'secret'),
            'charset'  => 'utf8',
            'prefix'   => '',
            'schema'   => 'public'
        ],
        'things-db' => [
            'driver'   => 'pgsql',
            'host'     => env('DB_THINGS_HOST', 'database_p'),
            'database' => env('DB_THINGS_DATABASE', 'dockerApp'), // This seems to be ignored
            'port'     => env('DB_THINGS_PORT', 5432),
            'username' => env('DB_THINGS_USERNAME', 'postgres'),
            'password' => env('DB_THINGS_PASSWORD', 'secret'),
            'charset'  => 'utf8',
            'prefix'   => '',
            'schema'   => 'public'
        ]
    ],
    'redis' => [
        'client' => env('REDIS_CLIENT', 'phpredis'),

        'cluster' => env('REDIS_CLUSTER', false),

        'default' => [
            'host' => env('REDIS_HOST', '127.0.0.1'),
            'password' => env('REDIS_PASSWORD', null),
            'port' => env('REDIS_PORT', 6379),
            'database' => env('REDIS_DB', 0),
        ],

        'cache' => [
            'host' => env('REDIS_HOST', '127.0.0.1'),
            'password' => env('REDIS_PASSWORD', null),
            'port' => env('REDIS_PORT', 6379),
            'database' => env('REDIS_CACHE_DB', 1),
        ],
    ],
];
