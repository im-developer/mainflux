<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

class CreateDevicesTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        Schema::connection('things-db')
            ->create('devices', function (Blueprint $table) {
            $table->uuid('id')->unique()->primary();
            $table->uuid('thing_id');
            $table->uuid('channel_id');
            $table->string('name');
            $table->string('icon')->nullable();
            $table->jsonb('metadata');
            $table->jsonb('pins');
            $table->boolean('active')->default(false);

//            $table->foreign('thing_id')
//                ->references('id')->on('things')
//                ->onDelete('cascade');

//            $table->foreign('channel_id')
//                ->references('id')->on('channels')
//                ->onDelete('cascade');
        });
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('devices');
    }
}
