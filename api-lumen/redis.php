<?php
require_once "vendor/autoload.php";

$eventLoop = \React\EventLoop\Factory::create();

$chat = new class implements \Laravie\Streaming\Listener {
    /**
     * @return array
     */
    public function subscribedChannels(): array {
        return ['mainflux.things'];
    }

    /**
     * @param  \Predis\Async\Client  $redis
     * @return void
     */
    public function onConnected($redis) {
        echo "Connected to redis!";
    }

    /**
     * @param  \Predis\Async\Client  $redis
     * @return void
     */
    public function onSubscribed($redis) {
        echo "Subscribed to channel `topic:*`!";
    }

    /**
     * Trigger on emitted listener.
     *
     * @param  object  $event
     * @param  object  $pubsub
     *
     * @return void
     */
    public function onEmitted($event, $pubsub) {
        // PUBLISH topic:laravel "Hello world"
echo "HI";
        # DESCRIBE $event
        #
        # {
        #   "kind": "pmessage",
        #   "pattern": "topic:*",
        #   "channel": "topic:laravel",
        #   "payload": "Hello world"
        # }
    }
};

$client = new Laravie\Streaming\Client(
    ['host' => '127.0.0.1', 'port' => 6378], $eventLoop
);

$client->connect($chat);

$eventLoop->run();
