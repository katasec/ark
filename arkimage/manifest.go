package arkimage

type Manifest struct {
	SchemaVersion int64               `json:"schemaVersion"`
	MediaType     string              `json:"mediaType"`
	Config        Config              `json:"config"`
	Layers        []Config            `json:"layers"`
	Annotations   ManifestAnnotations `json:"annotations"`
}

type ManifestAnnotations struct {
	ImageCreated     string `json:"org.opencontainers.image.created"`
	ImageDescription string `json:"org.opencontainers.image.description"`
	ImageType        string `json:"org.opencontainers.image.type"`
}

type Config struct {
	MediaType   string             `json:"mediaType"`
	Digest      string             `json:"digest"`
	Size        int64              `json:"size"`
	Annotations *ConfigAnnotations `json:"annotations,omitempty"`
}

type ConfigAnnotations struct {
	OrgOpencontainersImageTitle string `json:"org.opencontainers.image.title"`
}
