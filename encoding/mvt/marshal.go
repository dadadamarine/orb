package mvt

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/dadadamarine/orb"
	"strconv"

	"github.com/dadadamarine/orb/encoding/mvt/vectortile"
	"github.com/dadadamarine/orb/geojson"

	"github.com/gogo/protobuf/proto"
)

// MarshalGzipped will marshal the layers into Mapbox Vector Tile format
// and gzip the result. A lot of times MVT data is gzipped at rest,
// e.g. in a mbtiles file.
func MarshalGzipped(layers Layers) ([]byte, error) {
	data, err := Marshal(layers)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(nil)
	gzwriter := gzip.NewWriter(buf)

	_, err = gzwriter.Write(data)
	if err != nil {
		return nil, err
	}

	err = gzwriter.Close()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Marshal will take a set of layers and encode them into a Mapbox Vector Tile format.
func Marshal(layers Layers) ([]byte, error) {
	vt := &vectortile.Tile{
		Layers: make([]*vectortile.Tile_Layer, 0, len(layers)),
	}

	for _, l := range layers {
		v := l.Version
		e := l.Extent
		layer := &vectortile.Tile_Layer{
			Name:     &l.Name,
			Version:  &v,
			Extent:   &e,
			Features: make([]*vectortile.Tile_Feature, 0, len(l.Features)),
		}

		kve := newKeyValueEncoder()
		for _, f := range l.Features {
			if f.Geometry.GeoJSONType() == "GeometryCollection" {
				collection := f.Geometry.(orb.Collection)
				for _, g := range collection {
					err := addFeature(g, f.Properties, f.ID, kve, layer)
					if err != nil {
						continue
					}
				}
			}
			err := addFeature(f.Geometry, f.Properties, f.ID, kve, layer)
			if err != nil {
				continue
			}
		}

		layer.Keys = kve.Keys
		layer.Values = kve.Values

		vt.Layers = append(vt.Layers, layer)
	}

	return proto.Marshal(vt)
}

func addFeature(g orb.Geometry, p geojson.Properties, id interface{}, kve *keyValueEncoder, layer *vectortile.Tile_Layer) error {
	geomType, encodedGeometry, err := encodeGeometry(g)
	if err != nil {
		return err
	}
	if err != nil {
		return fmt.Errorf("error encoding geometry: %v", g)
	}

	tags, err := encodeProperties(kve, p)
	if err != nil {
		return fmt.Errorf("error encoding geometry: %v", g)
	}

	layer.Features = append(layer.Features, &vectortile.Tile_Feature{
		Id:       convertID(id),
		Tags:     tags,
		Type:     &geomType,
		Geometry: encodedGeometry,
	})
	return nil
}

func encodeProperties(kve *keyValueEncoder, properties geojson.Properties) ([]uint32, error) {
	tags := make([]uint32, 0, 2*len(properties))
	for k, v := range properties {
		ki := kve.Key(k)
		vi, err := kve.Value(v)
		if err != nil {
			return nil, fmt.Errorf("property %s: %v", k, err)
		}

		tags = append(tags, ki, vi)
	}

	return tags, nil
}

func convertID(id interface{}) *uint64 {
	if id == nil {
		return nil
	}

	switch id := id.(type) {
	case int:
		return convertIntID(id)
	case int8:
		return convertIntID(int(id))
	case int16:
		return convertIntID(int(id))
	case int32:
		return convertIntID(int(id))
	case int64:
		return convertIntID(int(id))
	case uint:
		v := uint64(id)
		return &v
	case uint8:
		v := uint64(id)
		return &v
	case uint16:
		v := uint64(id)
		return &v
	case uint32:
		v := uint64(id)
		return &v
	case uint64:
		v := uint64(id)
		return &v
	case float32:
		return convertIntID(int(id))
	case float64:
		return convertIntID(int(id))
	case string:
		i, err := strconv.Atoi(id)
		if err == nil {
			return convertIntID(i)
		}
	}

	return nil
}

func convertIntID(i int) *uint64 {
	if i < 0 {
		return nil
	}

	v := uint64(i)
	return &v
}
