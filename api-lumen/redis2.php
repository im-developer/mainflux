<?php
require_once "vendor/autoload.php";

use Asiries335\redisSteamPhp\Dto\StreamCommandCallTransporter;

// Example use in your app.
class Config implements \Asiries335\redisSteamPhp\ClientRedisStreamPhpInterface {

    private $client;

    public function __construct()
    {
        $this->client = new \Redis();
        $this->client->connect('127.0.0.1', '6378');
    }

    /**
     * Method for run command of redis
     *
     * @param StreamCommandCallTransporter $commandCallTransporter
     *
     * @return mixed
     *
     * @throws \Dto\Exceptions\InvalidDataTypeException
     * @throws \Dto\Exceptions\InvalidKeyException
     */
    public function call(StreamCommandCallTransporter $commandCallTransporter)
    {
        // Example use.
        return $this->client->rawCommand(
            $commandCallTransporter->get('command')->toScalar(),
            ...$commandCallTransporter->get('args')->toArray()
        );
    }
}

$client = new \Asiries335\redisSteamPhp\Client(new Config());
//$client->stream('test')->add(
//    'key',
//    [
//        'id'   => 123,
//        'name' => 'Barney',
//        'age'  => 25,
//    ]
//);
//
//$collection = $client->stream('mainflux.mqtt')->get();
//dd($collection);
$client->stream('mainflux.mqtt')->listen(
    function (\Asiries335\redisSteamPhp\Data\Message $message) {
        var_dump($message);
        // Your code...
    }
);
