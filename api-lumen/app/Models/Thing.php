<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Model;

class Thing extends Model
{
    /**
     * The connection name for the model.
     *
     * @var string
     */
    protected $connection = 'things-db';
    protected $keyType = 'string';
    protected $hidden = ['pivot'];

    public $timestamps = false;

    protected $casts = [
        'metadata' => 'object',
    ];

    public function channels() {
        return $this->belongsToMany(Channel::class, "connections");
    }

    public function devices() {
        return $this->hasMany(Device::class);
    }

}
