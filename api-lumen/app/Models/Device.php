<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Model;

class Device extends Model
{
    /**
     * The connection name for the model.
     *
     * @var string
     */
    protected $connection = 'things-db';
    protected $keyType = 'string';

    public $timestamps = false;

    protected $casts = [
        'metadata'  => 'object',
        'pins'      => 'array',
    ];

    public function thing() {
        return $this->belongsTo(Thing::class);
    }

    public function channel() {
        return $this->belongsTo(Channel::class);
    }

}
