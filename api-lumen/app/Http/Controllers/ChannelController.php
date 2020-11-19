<?php

namespace App\Http\Controllers;

use App\Models\Channel;
use Illuminate\Http\Request;
use Illuminate\Http\Response;
use Illuminate\Support\Str;

class ChannelController extends Controller
{
    /**
     * Create a new controller instance.
     *
     * @return void
     */
    public function __construct()
    {
        //
    }

    public function index() {
        $paginate = Channel::select('id', 'name', 'metadata')->paginate();

        return [
            'page'          => $paginate->currentPage(),
            'totalPages'    => $paginate->lastPage(),
            'totalRecords'  => $paginate->total(),
            'pageSize'      => $paginate->perPage(),
            'data'          => $paginate->items(),
        ];
    }

    public function all() {
        $query = Channel::select('id', 'name', 'metadata');

        if (\request()->has('things')) {
            $query->with('things:id,name,key,metadata');
        }

        if (\request()->has('devices')) {
            #$query->with('things:id,name,key,metadata');
            $query->with('devices');
        }

        return $query->get();
    }

    public function show($id) {
        return Channel::where('id', $id)->select('id', 'name', 'metadata')->firstOrFail();
    }

    public function things($id) {
        return Channel::with('things:id,name,key,metadata')->where('id', $id)->select('id', 'name', 'metadata')->firstOrFail();
    }

    public function devices($id) {
        return Channel::with('devices')->where('id', $id)->select('id', 'name', 'metadata')->firstOrFail();
    }

    public function store(Request $request) {
        $this->validate($request, [
            'name'  => 'required',
        ]);

        $name       = $request->get('name');
        $metadata   = $request->except('name');

        $channel = new Channel();
        $channel->id = Str::uuid();
        $channel->name = $name;
        $channel->owner = owner_id();
        $channel->metadata = $metadata;
        $channel->save();

        return response()->json($channel, Response::HTTP_CREATED);
    }

    public function update($id, Request $request) {
        $channel = Channel::where([
            'id'    => $id,
            'owner' => owner_id(),
        ])->firstOrFail();

        $this->validate($request, [
            'name'  => 'required',
        ]);

        $name       = $request->get('name');
        $metadata   = $request->except('name');

        $channel->name = $name;
        $channel->metadata = $metadata;
        $channel->save();

        return response()->json($channel, Response::HTTP_OK);
    }

    public function destroy($id) {
        $channel = Channel::where([
            'id'    => $id,
            'owner' => owner_id(),
        ])->firstOrFail();

        $channel->delete();

        return response()->json(null, Response::HTTP_NO_CONTENT);
    }
}
