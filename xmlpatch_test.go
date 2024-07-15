package xmlpatch

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestPatch(t *testing.T) {
	tests := []struct {
		name                string
		docDataFilepath     string
		xmlDiffDataFilepath string
		options             []Ops
		expectedFilepath    string
		wantErr             bool
	}{
		{
			name:                "replace existing attribute 'url' value of element 'project->component' identified by attribute 'name'",
			docDataFilepath:     "testdata/replace/attribute/1/workspace.before.xml",
			xmlDiffDataFilepath: "testdata/replace/attribute/1/diff.xml",
			expectedFilepath:    "testdata/replace/attribute/1/workspace.after.xml",
			wantErr:             false,
		},
		{
			name:                "set missing attribute 'url' value of element 'project->component' identified by attribute 'name'",
			docDataFilepath:     "testdata/replace/attribute/2/workspace.before.xml",
			xmlDiffDataFilepath: "testdata/replace/attribute/2/diff.xml",
			expectedFilepath:    "testdata/replace/attribute/2/workspace.after.xml",
			wantErr:             false,
		},
		{
			name:                "set attribute 'url' value of missing element 'project->component' identified by attribute 'name' with auto create on",
			docDataFilepath:     "testdata/replace/attribute/3/workspace.before.xml",
			xmlDiffDataFilepath: "testdata/replace/attribute/3/diff.xml",
			expectedFilepath:    "testdata/replace/attribute/3/workspace.after.xml",
			options:             []Ops{ReplaceAutoCreateMissing},
			wantErr:             false,
		},
		{
			name:                "set attribute 'url' value of missing element 'project->component' identified by attribute 'name' with auto create off",
			docDataFilepath:     "testdata/replace/attribute/3/workspace.before.xml",
			xmlDiffDataFilepath: "testdata/replace/attribute/3/diff.xml",
			expectedFilepath:    "testdata/replace/attribute/3/workspace.after.xml",
			wantErr:             true,
		},
		{
			name:                "set attribute 'url' value of missing element 'project->component->inner' identified by attribute 'name' with auto create on",
			docDataFilepath:     "testdata/replace/attribute/4/workspace.before.xml",
			xmlDiffDataFilepath: "testdata/replace/attribute/4/diff.xml",
			expectedFilepath:    "testdata/replace/attribute/4/workspace.after.xml",
			options:             []Ops{ReplaceAutoCreateMissing},
			wantErr:             false,
		},
		{
			name:                "set attribute 'url' value of 'project->component' identified by attribute 'name' with auto create on and doc data is empty",
			docDataFilepath:     "testdata/replace/attribute/5/workspace.before.xml",
			xmlDiffDataFilepath: "testdata/replace/attribute/5/diff.xml",
			expectedFilepath:    "testdata/replace/attribute/5/workspace.after.xml",
			options:             []Ops{ReplaceAutoCreateMissing},
			wantErr:             false,
		},
		{
			name:                "malformed target file",
			docDataFilepath:     "xmlpatch_test.go",
			xmlDiffDataFilepath: "testdata/replace/attribute/3/diff.xml",
			expectedFilepath:    "testdata/replace/attribute/3/workspace.after.xml",
			wantErr:             true,
		},
		{
			name:                "malformed diff file",
			docDataFilepath:     "testdata/replace/attribute/3/workspace.before.xml",
			xmlDiffDataFilepath: "xmlpatch_test.go",
			expectedFilepath:    "testdata/replace/attribute/3/workspace.after.xml",
			wantErr:             true,
		},
		{
			name:                "replace element",
			docDataFilepath:     "testdata/replace/element/1/domain.before.xml",
			xmlDiffDataFilepath: "testdata/replace/element/1/diff.xml",
			expectedFilepath:    "testdata/replace/element/1/domain.after.xml",
			wantErr:             false,
		},
		{
			name:                "add element",
			docDataFilepath:     "testdata/add/element/1/domain.before.xml",
			xmlDiffDataFilepath: "testdata/add/element/1/diff.xml",
			expectedFilepath:    "testdata/add/element/1/domain.after.xml",
			wantErr:             false,
		},
		{
			name:                "add duplicated element",
			docDataFilepath:     "testdata/add/element/2/domain.before.xml",
			xmlDiffDataFilepath: "testdata/add/element/2/diff.xml",
			expectedFilepath:    "testdata/add/element/2/domain.after.xml",
			wantErr:             false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			req := require.New(t)
			xmlDiffData, err := os.ReadFile(tt.xmlDiffDataFilepath)
			req.NoError(err)
			docData, err := os.ReadFile(tt.docDataFilepath)
			req.NoError(err)
			expected, err := os.ReadFile(tt.expectedFilepath)
			req.NoError(err)
			// when
			actual, err := Patch(docData, xmlDiffData, tt.options...)
			// then
			if tt.wantErr {
				req.Error(err)
				req.Nil(actual)
			} else {
				req.NoError(err)
				req.Equal(string(expected), string(actual))
			}
		})
	}
}
