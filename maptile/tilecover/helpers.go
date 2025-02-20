// Package tilecover computes the covering set of tiles for an orb.Geometry.
package tilecover

import (
	"fmt"

	"github.com/dadadamarine/orb"
	"github.com/dadadamarine/orb/maptile"
)

// Geometry returns the covering set of tiles for the given geometry.
func Geometry(g orb.Geometry, z maptile.Zoom) maptile.Set {
	if g == nil {
		return nil
	}

	switch g := g.(type) {
	case orb.Point:
		return Point(g, z)
	case orb.MultiPoint:
		return MultiPoint(g, z)
	case orb.LineString:
		return LineString(g, z)
	case orb.MultiLineString:
		return MultiLineString(g, z)
	case orb.Ring:
		return Ring(g, z)
	case orb.Polygon:
		return Polygon(g, z)
	case orb.MultiPolygon:
		return MultiPolygon(g, z)
	case orb.Collection:
		return Collection(g, z)
	case orb.Bound:
		return Bound(g, z)
	}

	panic(fmt.Sprintf("geometry type not supported: %T", g))
}

// Point creates a tile cover for the point, i.e. just the tile
// containing the point.
func Point(ll orb.Point, z maptile.Zoom) maptile.Set {
	return maptile.Set{
		maptile.At(ll, z): true,
	}
}

// MultiPoint creates a tile cover for the set of points,
func MultiPoint(mp orb.MultiPoint, z maptile.Zoom) maptile.Set {
	set := make(maptile.Set)
	for _, p := range mp {
		set[maptile.At(p, z)] = true
	}

	return set
}

// Bound creates a tile cover for the bound. i.e. all the tiles
// that intersect the bound.
func Bound(b orb.Bound, z maptile.Zoom) maptile.Set {
	lo := maptile.At(b.Min, z)
	hi := maptile.At(b.Max, z)

	result := make(maptile.Set, (hi.X-lo.X+1)*(lo.Y-hi.Y+1))

	for x := lo.X; x <= hi.X; x++ {
		for y := hi.Y; y <= lo.Y; y++ {
			result[maptile.Tile{X: x, Y: y, Z: z}] = true
		}
	}

	return result
}

// Collection returns the covering set of tiles for the
// geoemtry collection.
func Collection(c orb.Collection, z maptile.Zoom) maptile.Set {
	set := make(maptile.Set)
	for _, g := range c {
		set.Merge(Geometry(g, z))
	}

	return set
}
