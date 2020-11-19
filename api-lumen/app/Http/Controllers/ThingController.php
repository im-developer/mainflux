<?php

namespace App\Http\Controllers;

use App\Models\Thing;
use Illuminate\Http\Request;
use Illuminate\Http\Response;
use Illuminate\Support\Str;

class ThingController extends Controller
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
        $paginate = Thing::select('id', 'name', 'key', 'metadata')->paginate();

        return [
            'page' => $paginate->currentPage(),
            'totalPages' => $paginate->lastPage(),
            'totalRecords' => $paginate->total(),
            'pageSize' => $paginate->perPage(),
            'data' => $paginate->items(),
        ];
    }

    public function show($id) {
        return Thing::where('id', $id)->select('id', 'name', 'key', 'metadata')->firstOrFail();
    }

    public function store(Request $request) {
        $this->validate($request, [
            'name'      => 'required',
            'channel'   => 'required|exists:channels,id',
        ]);

        $name       = $request->get('name');
        $metadata   = $request->except(['name', 'channel']);

        $thing = new Thing();
        $thing->id = Str::uuid();
        $thing->key = Str::uuid();
        $thing->name = $name;
        $thing->owner = owner_id();
        $thing->metadata = $metadata;
        $thing->save();

        $thing->channels()->attach($request->get('channel'), [
            'thing_owner'   => owner_id(),
            'channel_owner' => owner_id(),
        ]);

        return response()->json($thing, Response::HTTP_CREATED);
    }

    public function update($id, Request $request) {
        $thing = Thing::where([
            'id'    => $id,
            'owner' => owner_id(),
        ])->firstOrFail();

        $this->validate($request, [
            'name'  => 'required',
        ]);

        $name       = $request->get('name');
        $metadata   = $request->except('name');

        $thing->name = $name;
        $thing->metadata = $metadata;
        $thing->save();

        return response()->json($thing, Response::HTTP_OK);
    }

    public function destroy($id) {
        $thing = Thing::where([
            'id'    => $id,
            'owner' => owner_id(),
        ])->firstOrFail();

        $thing->delete();

        return response()->json(null, Response::HTTP_NO_CONTENT);
    }
}
