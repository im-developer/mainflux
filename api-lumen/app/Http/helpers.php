<?php

function owner_id() {
    $user = request()->user();
    return $user->owner_id ?: $user->id;
}

function owner() {
    static $owner;
    if ($owner) return $owner;
    $user = request()->user();
    return $owner = $user->owner_id ? $user->owner : $user;
}
