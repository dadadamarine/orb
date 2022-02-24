package wkb

import (
	"testing"

	"github.com/dadadamarine/orb"
)

var (
	testPolygon = orb.Polygon{{
		{30, 10}, {40, 40}, {20, 40}, {10, 20}, {30, 10},
	}}
	testPolygonData = []byte{
		//01    02    03    04    05    06    07    08
		0x01, 0x03, 0x00, 0x00, 0x00,
		0x01, 0x00, 0x00, 0x00, // Number of Rings 1
		0x05, 0x00, 0x00, 0x00, // Number of Points 5
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x40, // X1 30
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x24, 0x40, // Y1 10
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x44, 0x40, // X2 40
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x44, 0x40, // Y2 40
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x34, 0x40, // X3 20
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x44, 0x40, // Y3 40
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x24, 0x40, // X4 10
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x34, 0x40, // Y4 20
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x40, // X5 30
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x24, 0x40, // Y5 10
	}
)

func TestPolygon(t *testing.T) {
	large := orb.Polygon{}
	for i := 0; i < maxMultiAlloc+100; i++ {
		large = append(large, orb.Ring{})
	}

	cases := []struct {
		name     string
		data     []byte
		expected orb.Polygon
	}{
		{
			name:     "polygon",
			data:     testPolygonData,
			expected: testPolygon,
		},
		{
			name:     "large",
			data:     MustMarshal(large),
			expected: large,
		},
		{
			name: "two ring polygon",
			data: []byte{
				//01    02    03    04    05    06    07    08
				0x01, 0x03, 0x00, 0x00, 0x00,
				0x02, 0x00, 0x00, 0x00, // Number of Lines (2)
				0x05, 0x00, 0x00, 0x00, // Number of Points in Line1 (5)
				0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x41, 0x40, // X1 35
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x24, 0x40, // Y1 10
				0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x46, 0x40, // X2 45
				0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x46, 0x40, // Y2 45
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x2e, 0x40, // X3 15
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x44, 0x40, // Y3 40
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x24, 0x40, // X4 10
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x34, 0x40, // Y4 20
				0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x41, 0x40, // X5 35
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x24, 0x40, // Y5 10
				0x04, 0x00, 0x00, 0x00, // Number of Points in Line2 (4)
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x34, 0x40, // X1 20
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x40, // Y1 30
				0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x41, 0x40, // X2 35
				0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x41, 0x40, // Y2 35
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x40, // X3 30
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x34, 0x40, // Y3 20
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x34, 0x40, // X4 20
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x40, // Y4 30
			},
			expected: orb.Polygon{
				{{35, 10}, {45, 45}, {15, 40}, {10, 20}, {35, 10}},
				{{20, 30}, {35, 35}, {30, 20}, {20, 30}},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			compare(t, tc.expected, tc.data)
		})
	}
}

var (
	testMultiPolygon = orb.MultiPolygon{
		{{{30, 20}, {45, 40}, {10, 40}, {30, 20}}},
		{{{15, 5}, {40, 10}, {10, 20}, {5, 10}, {15, 5}}},
	}
	testMultiPolygonData = []byte{
		//01    02    03    04    05    06    07    08
		0x01, 0x06, 0x00, 0x00, 0x00,
		0x02, 0x00, 0x00, 0x00, // Number of Polygons (2)
		0x01,                   // Byte Encoding Little
		0x03, 0x00, 0x00, 0x00, // Type Polygon1 (3)
		0x01, 0x00, 0x00, 0x00, // Number of Lines (1)
		0x04, 0x00, 0x00, 0x00, // Number of Points (4)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x40, // X1 30
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x34, 0x40, // Y1 20
		0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x46, 0x40, // X2 45
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x44, 0x40, // Y2 40
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x24, 0x40, // X3 10
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x44, 0x40, // Y3 40
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x40, // X4 30
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x34, 0x40, // X4 20
		0x01,                   // Byte Encoding Little
		0x03, 0x00, 0x00, 0x00, // Type Polygon2 (3)
		0x01, 0x00, 0x00, 0x00, // Number of Lines (1)
		0x05, 0x00, 0x00, 0x00, // Number of Points (5)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x2e, 0x40, // X1 15
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x14, 0x40, // Y1  5
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x44, 0x40, // X2 40
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x24, 0x40, // Y2 10
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x24, 0x40, // X3 10
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x34, 0x40, // Y3 20
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x14, 0x40, // X4  5
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x24, 0x40, // Y4 10
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x2e, 0x40, // X5 15
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x14, 0x40, // Y5  5
	}

	testMultiPolygonSingle = orb.MultiPolygon{
		{
			{{20, 35}, {10, 30}, {10, 10}, {30, 5}, {45, 20}, {20, 35}},
			{{30, 20}, {20, 15}, {20, 25}, {30, 20}}},
	}
	testMultiPolygonSingleData = []byte{
		//01    02    03    04    05    06    07    08
		0x01, 0x06, 0x00, 0x00, 0x00,
		0x01, 0x00, 0x00, 0x00, // Number of Polygons (1)
		0x01,                   // Byte order marker little
		0x03, 0x00, 0x00, 0x00, // Type Polygon(3)
		0x02, 0x00, 0x00, 0x00, // Number of Lines(2)
		0x06, 0x00, 0x00, 0x00, // Number of Points(6)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x34, 0x40, // X1 20
		0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x41, 0x40, // Y1 35
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x24, 0x40, // X2 10
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x40, // Y2 30
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x24, 0x40, // X3 10
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x24, 0x40, // Y3 10
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x40, // X4 30
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x14, 0x40, // Y4 5
		0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x46, 0x40, // X5 45
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x34, 0x40, // Y5 20
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x34, 0x40, // X6 20
		0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x41, 0x40, // Y6 35
		0x04, 0x00, 0x00, 0x00, // Number of Points(4)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x40, // X1 30
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x34, 0x40, // Y1 20
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x34, 0x40, // X2 20
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x2e, 0x40, // Y2 15
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x34, 0x40, // X3 20
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x39, 0x40, // Y3 25
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x40, // X4 30
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x34, 0x40, // Y4 20
	}
)

func TestMultiPolygon(t *testing.T) {
	large := orb.MultiPolygon{}
	for i := 0; i < maxMultiAlloc+100; i++ {
		large = append(large, orb.Polygon{})
	}

	cases := []struct {
		name     string
		data     []byte
		expected orb.MultiPolygon
	}{
		{
			name:     "multi polygon",
			data:     testMultiPolygonData,
			expected: testMultiPolygon,
		},
		{
			name:     "single multi polygon",
			data:     testMultiPolygonSingleData,
			expected: testMultiPolygonSingle,
		},
		{
			name:     "large",
			data:     MustMarshal(large),
			expected: large,
		},
		{
			name: "three polygons",
			data: []byte{
				//01    02    03    04    05    06    07    08
				0x01, 0x06, 0x00, 0x00, 0x00,
				0x02, 0x00, 0x00, 0x00, // Number of Polygons (2)
				0x01,                   // Byte order marker little
				0x03, 0x00, 0x00, 0x00, // type Polygon (3)
				0x01, 0x00, 0x00, 0x00, // Number of Lines (1)
				0x04, 0x00, 0x00, 0x00, // Number of Points (4)
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x44, 0x40, // X1 40
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x44, 0x40, // Y1 40
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x34, 0x40, // X2 20
				0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x46, 0x40, // Y2 45
				0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x46, 0x40, // X3 45
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x40, // Y3 30
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x44, 0x40, // X4 40
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x44, 0x40, // Y4 40
				0x01,                   // Byte order marker little
				0x03, 0x00, 0x00, 0x00, // Type Polygon(3)
				0x02, 0x00, 0x00, 0x00, // Number of Lines(2)
				0x06, 0x00, 0x00, 0x00, // Number of Points(6)
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x34, 0x40, // X1 20
				0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x41, 0x40, // Y1 35
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x24, 0x40, // X2 10
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x40, // Y2 30
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x24, 0x40, // X3 10
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x24, 0x40, // Y3 10
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x40, // X4 30
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x14, 0x40, // Y4 5
				0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x46, 0x40, // X5 45
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x34, 0x40, // Y5 20
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x34, 0x40, // X6 20
				0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x41, 0x40, // Y6 35
				0x04, 0x00, 0x00, 0x00, // Number of Points(4)
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x40, // X1 30
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x34, 0x40, // Y1 20
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x34, 0x40, // X2 20
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x2e, 0x40, // Y2 15
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x34, 0x40, // X3 20
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x39, 0x40, // Y3 25
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x40, // X4 30
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x34, 0x40, // Y4 20

			},
			expected: orb.MultiPolygon{
				{
					{{40, 40}, {20, 45}, {45, 30}, {40, 40}},
				},
				{
					{{20, 35}, {10, 30}, {10, 10}, {30, 5}, {45, 20}, {20, 35}},
					{{30, 20}, {20, 15}, {20, 25}, {30, 20}},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			compare(t, tc.expected, tc.data)
		})
	}
}
