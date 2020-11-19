<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Model;

class Channel extends Model
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
        'metadata' => 'object',
    ];
    protected $hidden = ['pivot'];

    public function things() {
        return $this->belongsToMany(Thing::class, "connections");
    }

    public function devices() {
        return $this->hasMany(Device::class);
    }
}
