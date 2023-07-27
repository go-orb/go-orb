package metadata

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMetadataSet(t *testing.T) {
	ctx := Set(context.TODO(), "Key", "val")

	val, ok := Get(ctx, "Key")
	if !ok {
		t.Fatal("key Key not found")
	}
	if val != "val" {
		t.Errorf("key Key with value val != %v", val)
	}
}

func TestMetadataDelete(t *testing.T) {
	md := Metadata{
		"foo": "bar",
		"baz": "empty",
	}

	ctx := md.To(context.TODO())
	ctx = Delete(ctx, "Baz")

	emd, ok := From(ctx)
	if !ok {
		t.Fatal("key Key not found")
	}

	_, ok = emd["Baz"]
	if ok {
		t.Fatal("key Baz not deleted")
	}
}

func TestMetadataCopy(t *testing.T) {
	md := Metadata{
		"foo": "bar",
		"bar": "baz",
	}

	cp := Copy(md)

	for k, v := range md {
		if cv := cp[k]; cv != v {
			t.Fatalf("Got %s:%s for %s:%s", k, cv, k, v)
		}
	}
}

func TestMetadataContext(t *testing.T) {
	md := Metadata{
		"foo": "bar",
	}

	ctx := md.To(context.TODO())

	emd, ok := From(ctx)
	if !ok {
		t.Errorf("Unexpected error retrieving metadata, got %t", ok)
	}

	if emd["foo"] != md["foo"] {
		t.Errorf("Expected key: %s val: %s, got key: %s val: %s", "Foo", md["Foo"], "Foo", emd["Foo"])
	}

	if i := len(emd); i != 1 {
		t.Errorf("Expected metadata length 1 got %d", i)
	}
}

func TestMergeContext(t *testing.T) {
	type args struct {
		existing  Metadata
		append    Metadata
		overwrite bool
	}
	tests := []struct {
		name string
		args args
		want Metadata
	}{
		{
			name: "matching key, overwrite false",
			args: args{
				existing:  Metadata{"foo": "bar", "sumo": "demo"},
				append:    Metadata{"sumo": "demo2"},
				overwrite: false,
			},
			want: Metadata{"foo": "bar", "sumo": "demo"},
		},
		{
			name: "matching key, overwrite true",
			args: args{
				existing:  Metadata{"foo": "bar", "sumo": "demo"},
				append:    Metadata{"sumo": "demo2"},
				overwrite: true,
			},
			want: Metadata{"foo": "bar", "sumo": "demo2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := From(Merge(tt.args.existing.To(context.TODO()), tt.args.append, tt.args.overwrite))
			require.Equal(t, got, tt.want)
		})
	}
}
