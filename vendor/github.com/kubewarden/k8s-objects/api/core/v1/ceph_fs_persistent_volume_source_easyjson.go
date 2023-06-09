// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package v1

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjsonB1662561DecodeGithubComKubewardenK8sObjectsApiCoreV1(in *jlexer.Lexer, out *CephFSPersistentVolumeSource) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "monitors":
			if in.IsNull() {
				in.Skip()
				out.Monitors = nil
			} else {
				in.Delim('[')
				if out.Monitors == nil {
					if !in.IsDelim(']') {
						out.Monitors = make([]string, 0, 4)
					} else {
						out.Monitors = []string{}
					}
				} else {
					out.Monitors = (out.Monitors)[:0]
				}
				for !in.IsDelim(']') {
					var v1 string
					v1 = string(in.String())
					out.Monitors = append(out.Monitors, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "path":
			out.Path = string(in.String())
		case "readOnly":
			out.ReadOnly = bool(in.Bool())
		case "secretFile":
			out.SecretFile = string(in.String())
		case "secretRef":
			if in.IsNull() {
				in.Skip()
				out.SecretRef = nil
			} else {
				if out.SecretRef == nil {
					out.SecretRef = new(SecretReference)
				}
				easyjsonB1662561DecodeGithubComKubewardenK8sObjectsApiCoreV11(in, out.SecretRef)
			}
		case "user":
			out.User = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonB1662561EncodeGithubComKubewardenK8sObjectsApiCoreV1(out *jwriter.Writer, in CephFSPersistentVolumeSource) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"monitors\":"
		out.RawString(prefix[1:])
		if in.Monitors == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v2, v3 := range in.Monitors {
				if v2 > 0 {
					out.RawByte(',')
				}
				out.String(string(v3))
			}
			out.RawByte(']')
		}
	}
	if in.Path != "" {
		const prefix string = ",\"path\":"
		out.RawString(prefix)
		out.String(string(in.Path))
	}
	if in.ReadOnly {
		const prefix string = ",\"readOnly\":"
		out.RawString(prefix)
		out.Bool(bool(in.ReadOnly))
	}
	if in.SecretFile != "" {
		const prefix string = ",\"secretFile\":"
		out.RawString(prefix)
		out.String(string(in.SecretFile))
	}
	if in.SecretRef != nil {
		const prefix string = ",\"secretRef\":"
		out.RawString(prefix)
		easyjsonB1662561EncodeGithubComKubewardenK8sObjectsApiCoreV11(out, *in.SecretRef)
	}
	if in.User != "" {
		const prefix string = ",\"user\":"
		out.RawString(prefix)
		out.String(string(in.User))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v CephFSPersistentVolumeSource) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonB1662561EncodeGithubComKubewardenK8sObjectsApiCoreV1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v CephFSPersistentVolumeSource) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonB1662561EncodeGithubComKubewardenK8sObjectsApiCoreV1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *CephFSPersistentVolumeSource) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonB1662561DecodeGithubComKubewardenK8sObjectsApiCoreV1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *CephFSPersistentVolumeSource) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonB1662561DecodeGithubComKubewardenK8sObjectsApiCoreV1(l, v)
}
func easyjsonB1662561DecodeGithubComKubewardenK8sObjectsApiCoreV11(in *jlexer.Lexer, out *SecretReference) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "name":
			out.Name = string(in.String())
		case "namespace":
			out.Namespace = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonB1662561EncodeGithubComKubewardenK8sObjectsApiCoreV11(out *jwriter.Writer, in SecretReference) {
	out.RawByte('{')
	first := true
	_ = first
	if in.Name != "" {
		const prefix string = ",\"name\":"
		first = false
		out.RawString(prefix[1:])
		out.String(string(in.Name))
	}
	if in.Namespace != "" {
		const prefix string = ",\"namespace\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Namespace))
	}
	out.RawByte('}')
}
