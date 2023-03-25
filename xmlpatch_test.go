package xmlpatch

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestPatch(t *testing.T) {
	req := require.New(t)
	diffXml, err := os.ReadFile("testdata/replace/attribute/1/diff.xml")
	req.NoError(err)
	workspaceBeforeXml, err := os.ReadFile("testdata/replace/attribute/1/workspace.before.xml")
	req.NoError(err)
	workspaceAfterXml, err := os.ReadFile("testdata/replace/attribute/1/workspace.after.xml")
	req.NoError(err)
	type args struct {
		docData     []byte
		xmlDiffData []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "replace attribute 'url' value of element 'project->component' identified by attribute 'name'",
			args: args{
				docData:     workspaceBeforeXml,
				xmlDiffData: diffXml,
			},
			want:    workspaceAfterXml,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Patch(tt.args.docData, tt.args.xmlDiffData)
			if tt.wantErr {
				req.Error(err)
			} else {
				req.NoError(err)
			}
			req.Equal(string(tt.want), string(got))
		})
	}
}
