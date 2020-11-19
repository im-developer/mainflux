<?php

/** @var \Laravel\Lumen\Routing\Router $router */

/*
|--------------------------------------------------------------------------
| Application Routes
|--------------------------------------------------------------------------
|
| Here is where you can register all of the routes for an application.
| It is a breeze. Simply tell Lumen the URIs it should respond to
| and give it the Closure to call when that URI is requested.
|
*/

$router->get('/', function () use ($router) {
    return $router->app->version();
});

# Groups as homes [ add by default with user ]
# Users [ add, edit, delete ]
# Things [ add, edit, delete ]
# Channels [ add, edit, delete ]

$router->group(['prefix' => 'v1'], function () use ($router) {
    $router->post('auth/login', 'AuthController@login');
    $router->post('logout', 'AuthController@logout');
    $router->post('auth/refresh', 'AuthController@refresh');
    $router->post('me', 'AuthController@me');

    $router->group(['middleware' => 'auth:api'], function () use ($router) {
        $router->get('channels', 'ChannelController@index');
        $router->get('channels/all', 'ChannelController@all');
        $router->post('channels', 'ChannelController@store');
        $router->get('channels/{id}', 'ChannelController@show');
        $router->put('channels/{id}', 'ChannelController@update');
        $router->delete('channels/{id}', 'ChannelController@destroy');
        $router->get('channels/{id}/things', 'ChannelController@things');
        $router->get('channels/{id}/devices', 'ChannelController@devices');

        $router->get('things', 'ThingController@index');
        $router->post('things', 'ThingController@store');
        $router->get('things/{id}', 'ThingController@show');
        $router->put('things/{id}', 'ThingController@update');
        $router->delete('things/{id}', 'ThingController@destroy');
    });
});


