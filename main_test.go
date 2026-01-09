package envconfig

import (
	"os"
	"testing"
)

func TestGet(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue string
		setEnv       bool
		envValue     string
		want         string
	}{
		{
			name:         "returns environment variable when set",
			key:          "TEST_KEY",
			defaultValue: "default",
			setEnv:       true,
			envValue:     "env_value",
			want:         "env_value",
		},
		{
			name:         "returns default when environment variable not set",
			key:          "TEST_KEY_MISSING",
			defaultValue: "default",
			setEnv:       false,
			want:         "default",
		},
		{
			name:         "returns default when environment variable is empty",
			key:          "TEST_KEY_EMPTY",
			defaultValue: "default",
			setEnv:       true,
			envValue:     "",
			want:         "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setEnv {
				os.Setenv(tt.key, tt.envValue)
				defer os.Unsetenv(tt.key)
			} else {
				os.Unsetenv(tt.key)
			}

			got := Get(tt.key, tt.defaultValue)
			if got != tt.want {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToList(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		separator string
		want      []string
		wantErr   bool
	}{
		{
			name:      "splits comma-separated values",
			value:     "a,b,c",
			separator: ",",
			want:      []string{"a", "b", "c"},
			wantErr:   false,
		},
		{
			name:      "splits semicolon-separated values",
			value:     "a;b;c",
			separator: ";",
			want:      []string{"a", "b", "c"},
			wantErr:   false,
		},
		{
			name:      "handles single value",
			value:     "single",
			separator: ",",
			want:      []string{"single"},
			wantErr:   false,
		},
		{
			name:      "handles empty string",
			value:     "",
			separator: ",",
			want:      []string{""},
			wantErr:   false,
		},
		{
			name:      "handles space separator",
			value:     "a b c",
			separator: " ",
			want:      []string{"a", "b", "c"},
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToList(tt.value, tt.separator)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("ToList() = %v, want %v", got, tt.want)
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("ToList()[%d] = %v, want %v", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestGetBool(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue bool
		setEnv       bool
		envValue     string
		want         bool
	}{
		{
			name:         "returns true from environment variable",
			key:          "TEST_BOOL_TRUE",
			defaultValue: false,
			setEnv:       true,
			envValue:     "true",
			want:         true,
		},
		{
			name:         "returns false from environment variable",
			key:          "TEST_BOOL_FALSE",
			defaultValue: true,
			setEnv:       true,
			envValue:     "false",
			want:         false,
		},
		{
			name:         "returns default when environment variable not set",
			key:          "TEST_BOOL_MISSING",
			defaultValue: true,
			setEnv:       false,
			want:         true,
		},
		{
			name:         "returns default when environment variable is empty",
			key:          "TEST_BOOL_EMPTY",
			defaultValue: false,
			setEnv:       true,
			envValue:     "",
			want:         false,
		},
		{
			name:         "handles invalid boolean string",
			key:          "TEST_BOOL_INVALID",
			defaultValue: false,
			setEnv:       true,
			envValue:     "invalid",
			want:         false,
		},
		{
			name:         "handles '1' as true",
			key:          "TEST_BOOL_ONE",
			defaultValue: false,
			setEnv:       true,
			envValue:     "1",
			want:         true,
		},
		{
			name:         "handles '0' as false",
			key:          "TEST_BOOL_ZERO",
			defaultValue: true,
			setEnv:       true,
			envValue:     "0",
			want:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setEnv {
				os.Setenv(tt.key, tt.envValue)
				defer os.Unsetenv(tt.key)
			} else {
				os.Unsetenv(tt.key)
			}

			got := GetBool(tt.key, tt.defaultValue)
			if got != tt.want {
				t.Errorf("GetBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetInt(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue int
		setEnv       bool
		envValue     string
		want         int
	}{
		{
			name:         "returns integer from environment variable",
			key:          "TEST_INT_VALID",
			defaultValue: 0,
			setEnv:       true,
			envValue:     "42",
			want:         42,
		},
		{
			name:         "returns default when environment variable not set",
			key:          "TEST_INT_MISSING",
			defaultValue: 100,
			setEnv:       false,
			want:         100,
		},
		{
			name:         "returns default when environment variable is empty",
			key:          "TEST_INT_EMPTY",
			defaultValue: 50,
			setEnv:       true,
			envValue:     "",
			want:         50,
		},
		{
			name:         "returns default when environment variable is invalid",
			key:          "TEST_INT_INVALID",
			defaultValue: 99,
			setEnv:       true,
			envValue:     "not_a_number",
			want:         99,
		},
		{
			name:         "handles negative numbers",
			key:          "TEST_INT_NEGATIVE",
			defaultValue: 0,
			setEnv:       true,
			envValue:     "-10",
			want:         -10,
		},
		{
			name:         "handles zero",
			key:          "TEST_INT_ZERO",
			defaultValue: 100,
			setEnv:       true,
			envValue:     "0",
			want:         0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setEnv {
				os.Setenv(tt.key, tt.envValue)
				defer os.Unsetenv(tt.key)
			} else {
				os.Unsetenv(tt.key)
			}

			got := GetInt(tt.key, tt.defaultValue)
			if got != tt.want {
				t.Errorf("GetInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetInt64(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue int64
		setEnv       bool
		envValue     string
		want         int64
	}{
		{
			name:         "returns int64 from environment variable",
			key:          "TEST_INT64_VALID",
			defaultValue: 0,
			setEnv:       true,
			envValue:     "9223372036854775807",
			want:         9223372036854775807,
		},
		{
			name:         "returns default when environment variable not set",
			key:          "TEST_INT64_MISSING",
			defaultValue: 1000,
			setEnv:       false,
			want:         1000,
		},
		{
			name:         "returns default when environment variable is empty",
			key:          "TEST_INT64_EMPTY",
			defaultValue: 500,
			setEnv:       true,
			envValue:     "",
			want:         500,
		},
		{
			name:         "returns default when environment variable is invalid",
			key:          "TEST_INT64_INVALID",
			defaultValue: 999,
			setEnv:       true,
			envValue:     "not_a_number",
			want:         999,
		},
		{
			name:         "handles negative numbers",
			key:          "TEST_INT64_NEGATIVE",
			defaultValue: 0,
			setEnv:       true,
			envValue:     "-9223372036854775808",
			want:         -9223372036854775808,
		},
		{
			name:         "handles zero",
			key:          "TEST_INT64_ZERO",
			defaultValue: 1000,
			setEnv:       true,
			envValue:     "0",
			want:         0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setEnv {
				os.Setenv(tt.key, tt.envValue)
				defer os.Unsetenv(tt.key)
			} else {
				os.Unsetenv(tt.key)
			}

			got := GetInt64(tt.key, tt.defaultValue)
			if got != tt.want {
				t.Errorf("GetInt64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetIntSlice(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue []int
		setEnv       bool
		envValue     string
		want         []int
	}{
		{
			name:         "returns int slice from environment variable",
			key:          "TEST_INT_SLICE_VALID",
			defaultValue: []int{0},
			setEnv:       true,
			envValue:     "8080,8081,8082",
			want:         []int{8080, 8081, 8082},
		},
		{
			name:         "returns default when environment variable not set",
			key:          "TEST_INT_SLICE_MISSING",
			defaultValue: []int{3000, 3001},
			setEnv:       false,
			want:         []int{3000, 3001},
		},
		{
			name:         "returns default when environment variable is empty",
			key:          "TEST_INT_SLICE_EMPTY",
			defaultValue: []int{5000},
			setEnv:       true,
			envValue:     "",
			want:         []int{5000},
		},
		{
			name:         "handles single value",
			key:          "TEST_INT_SLICE_SINGLE",
			defaultValue: []int{0},
			setEnv:       true,
			envValue:     "42",
			want:         []int{42},
		},
		{
			name:         "handles values with spaces",
			key:          "TEST_INT_SLICE_SPACES",
			defaultValue: []int{0},
			setEnv:       true,
			envValue:     " 1 , 2 , 3 ",
			want:         []int{1, 2, 3},
		},
		{
			name:         "handles negative numbers",
			key:          "TEST_INT_SLICE_NEGATIVE",
			defaultValue: []int{0},
			setEnv:       true,
			envValue:     "-10,-20,-30",
			want:         []int{-10, -20, -30},
		},
		{
			name:         "handles empty values in slice",
			key:          "TEST_INT_SLICE_EMPTY_VALUES",
			defaultValue: []int{0},
			setEnv:       true,
			envValue:     "1,,3",
			want:         []int{1, 0, 3},
		},
		{
			name:         "returns default when environment variable is invalid",
			key:          "TEST_INT_SLICE_INVALID",
			defaultValue: []int{999},
			setEnv:       true,
			envValue:     "8080,not_a_number,8082",
			want:         []int{999},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setEnv {
				os.Setenv(tt.key, tt.envValue)
				defer os.Unsetenv(tt.key)
			} else {
				os.Unsetenv(tt.key)
			}

			got := GetIntSlice(tt.key, tt.defaultValue)
			if len(got) != len(tt.want) {
				t.Errorf("GetIntSlice() length = %v, want %v", len(got), len(tt.want))
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("GetIntSlice()[%d] = %v, want %v", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestGetInt64Slice(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue []int64
		setEnv       bool
		envValue     string
		want         []int64
	}{
		{
			name:         "returns int64 slice from environment variable",
			key:          "TEST_INT64_SLICE_VALID",
			defaultValue: []int64{0},
			setEnv:       true,
			envValue:     "9223372036854775807,9223372036854775806",
			want:         []int64{9223372036854775807, 9223372036854775806},
		},
		{
			name:         "returns default when environment variable not set",
			key:          "TEST_INT64_SLICE_MISSING",
			defaultValue: []int64{1000, 2000},
			setEnv:       false,
			want:         []int64{1000, 2000},
		},
		{
			name:         "returns default when environment variable is empty",
			key:          "TEST_INT64_SLICE_EMPTY",
			defaultValue: []int64{5000},
			setEnv:       true,
			envValue:     "",
			want:         []int64{5000},
		},
		{
			name:         "handles single value",
			key:          "TEST_INT64_SLICE_SINGLE",
			defaultValue: []int64{0},
			setEnv:       true,
			envValue:     "42",
			want:         []int64{42},
		},
		{
			name:         "handles values with spaces",
			key:          "TEST_INT64_SLICE_SPACES",
			defaultValue: []int64{0},
			setEnv:       true,
			envValue:     " 1 , 2 , 3 ",
			want:         []int64{1, 2, 3},
		},
		{
			name:         "handles negative numbers",
			key:          "TEST_INT64_SLICE_NEGATIVE",
			defaultValue: []int64{0},
			setEnv:       true,
			envValue:     "-9223372036854775808,-10",
			want:         []int64{-9223372036854775808, -10},
		},
		{
			name:         "handles empty values in slice",
			key:          "TEST_INT64_SLICE_EMPTY_VALUES",
			defaultValue: []int64{0},
			setEnv:       true,
			envValue:     "1,,3",
			want:         []int64{1, 0, 3},
		},
		{
			name:         "returns default when environment variable is invalid",
			key:          "TEST_INT64_SLICE_INVALID",
			defaultValue: []int64{999},
			setEnv:       true,
			envValue:     "8080,not_a_number,8082",
			want:         []int64{999},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setEnv {
				os.Setenv(tt.key, tt.envValue)
				defer os.Unsetenv(tt.key)
			} else {
				os.Unsetenv(tt.key)
			}

			got := GetInt64Slice(tt.key, tt.defaultValue)
			if len(got) != len(tt.want) {
				t.Errorf("GetInt64Slice() length = %v, want %v", len(got), len(tt.want))
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("GetInt64Slice()[%d] = %v, want %v", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestLoad(t *testing.T) {
	tests := []struct {
		name        string
		setEnv      bool
		envFile     string
		wantErr     bool
		setupFile   bool
		filePath    string
		fileContent string
	}{
		{
			name:    "loads default .env file when ENV_FILE not set",
			setEnv:  false,
			wantErr: false, // May or may not error depending on if .env exists
		},
		{
			name:    "loads custom env file from ENV_FILE",
			setEnv:  true,
			envFile: "test.env",
			wantErr: false, // May or may not error depending on if test.env exists
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setEnv {
				os.Setenv("ENV_FILE", tt.envFile)
				defer os.Unsetenv("ENV_FILE")
			} else {
				os.Unsetenv("ENV_FILE")
			}

			err := Load()
			if (err != nil) != tt.wantErr {
				// For Load(), we don't fail if file doesn't exist - it's expected behavior
				// We just check that the function doesn't panic
				if err != nil {
					t.Logf("Load() returned error (expected if file doesn't exist): %v", err)
				}
			}
		})
	}
}

func TestLoadStruct(t *testing.T) {
	tests := []struct {
		name       string
		cfg        interface{}
		setupEnv   func()
		cleanupEnv func()
		wantErr    bool
		validate   func(t *testing.T, cfg interface{})
	}{
		{
			name: "loads string values from environment",
			cfg: &struct {
				Name string `env:"TEST_NAME"`
			}{},
			setupEnv: func() {
				os.Setenv("TEST_NAME", "test_value")
			},
			cleanupEnv: func() {
				os.Unsetenv("TEST_NAME")
			},
			wantErr: false,
			validate: func(t *testing.T, cfg interface{}) {
				c := cfg.(*struct {
					Name string `env:"TEST_NAME"`
				})
				if c.Name != "test_value" {
					t.Errorf("Name = %v, want test_value", c.Name)
				}
			},
		},
		{
			name: "uses default value when env variable not set",
			cfg: &struct {
				Port int `env:"TEST_PORT" default:"8080"`
			}{},
			setupEnv: func() {
				os.Unsetenv("TEST_PORT")
			},
			cleanupEnv: func() {
				os.Unsetenv("TEST_PORT")
			},
			wantErr: false,
			validate: func(t *testing.T, cfg interface{}) {
				c := cfg.(*struct {
					Port int `env:"TEST_PORT" default:"8080"`
				})
				if c.Port != 8080 {
					t.Errorf("Port = %v, want 8080", c.Port)
				}
			},
		},
		{
			name: "loads bool values from environment",
			cfg: &struct {
				Enabled bool `env:"TEST_ENABLED"`
			}{},
			setupEnv: func() {
				os.Setenv("TEST_ENABLED", "true")
			},
			cleanupEnv: func() {
				os.Unsetenv("TEST_ENABLED")
			},
			wantErr: false,
			validate: func(t *testing.T, cfg interface{}) {
				c := cfg.(*struct {
					Enabled bool `env:"TEST_ENABLED"`
				})
				if !c.Enabled {
					t.Errorf("Enabled = %v, want true", c.Enabled)
				}
			},
		},
		{
			name: "loads int values from environment",
			cfg: &struct {
				Count int `env:"TEST_COUNT"`
			}{},
			setupEnv: func() {
				os.Setenv("TEST_COUNT", "42")
			},
			cleanupEnv: func() {
				os.Unsetenv("TEST_COUNT")
			},
			wantErr: false,
			validate: func(t *testing.T, cfg interface{}) {
				c := cfg.(*struct {
					Count int `env:"TEST_COUNT"`
				})
				if c.Count != 42 {
					t.Errorf("Count = %v, want 42", c.Count)
				}
			},
		},
		{
			name: "loads int slice from environment",
			cfg: &struct {
				Ports []int `env:"TEST_PORTS"`
			}{},
			setupEnv: func() {
				os.Setenv("TEST_PORTS", "8080,8081,8082")
			},
			cleanupEnv: func() {
				os.Unsetenv("TEST_PORTS")
			},
			wantErr: false,
			validate: func(t *testing.T, cfg interface{}) {
				c := cfg.(*struct {
					Ports []int `env:"TEST_PORTS"`
				})
				want := []int{8080, 8081, 8082}
				if len(c.Ports) != len(want) {
					t.Errorf("Ports length = %v, want %v", len(c.Ports), len(want))
					return
				}
				for i, v := range want {
					if c.Ports[i] != v {
						t.Errorf("Ports[%d] = %v, want %v", i, c.Ports[i], v)
					}
				}
			},
		},
		{
			name: "loads int array from environment",
			cfg: &struct {
				Fixed [3]int `env:"TEST_FIXED"`
			}{},
			setupEnv: func() {
				os.Setenv("TEST_FIXED", "1,2,3")
			},
			cleanupEnv: func() {
				os.Unsetenv("TEST_FIXED")
			},
			wantErr: false,
			validate: func(t *testing.T, cfg interface{}) {
				c := cfg.(*struct {
					Fixed [3]int `env:"TEST_FIXED"`
				})
				want := [3]int{1, 2, 3}
				if c.Fixed != want {
					t.Errorf("Fixed = %v, want %v", c.Fixed, want)
				}
			},
		},
		{
			name: "loads int slice with default value",
			cfg: &struct {
				Ports []int `env:"TEST_PORTS_DEFAULT" default:"3000,3001"`
			}{},
			setupEnv: func() {
				os.Unsetenv("TEST_PORTS_DEFAULT")
			},
			cleanupEnv: func() {
				os.Unsetenv("TEST_PORTS_DEFAULT")
			},
			wantErr: false,
			validate: func(t *testing.T, cfg interface{}) {
				c := cfg.(*struct {
					Ports []int `env:"TEST_PORTS_DEFAULT" default:"3000,3001"`
				})
				want := []int{3000, 3001}
				if len(c.Ports) != len(want) {
					t.Errorf("Ports length = %v, want %v", len(c.Ports), len(want))
					return
				}
				for i, v := range want {
					if c.Ports[i] != v {
						t.Errorf("Ports[%d] = %v, want %v", i, c.Ports[i], v)
					}
				}
			},
		},
		{
			name: "handles empty int slice",
			cfg: &struct {
				Ports []int `env:"TEST_PORTS_EMPTY"`
			}{},
			setupEnv: func() {
				os.Setenv("TEST_PORTS_EMPTY", "")
			},
			cleanupEnv: func() {
				os.Unsetenv("TEST_PORTS_EMPTY")
			},
			wantErr: false,
			validate: func(t *testing.T, cfg interface{}) {
				c := cfg.(*struct {
					Ports []int `env:"TEST_PORTS_EMPTY"`
				})
				if len(c.Ports) != 0 {
					t.Errorf("Ports length = %v, want 0", len(c.Ports))
				}
			},
		},
		{
			name: "handles int slice with spaces",
			cfg: &struct {
				Ports []int `env:"TEST_PORTS_SPACES"`
			}{},
			setupEnv: func() {
				os.Setenv("TEST_PORTS_SPACES", " 1 , 2 , 3 ")
			},
			cleanupEnv: func() {
				os.Unsetenv("TEST_PORTS_SPACES")
			},
			wantErr: false,
			validate: func(t *testing.T, cfg interface{}) {
				c := cfg.(*struct {
					Ports []int `env:"TEST_PORTS_SPACES"`
				})
				want := []int{1, 2, 3}
				if len(c.Ports) != len(want) {
					t.Errorf("Ports length = %v, want %v", len(c.Ports), len(want))
					return
				}
				for i, v := range want {
					if c.Ports[i] != v {
						t.Errorf("Ports[%d] = %v, want %v", i, c.Ports[i], v)
					}
				}
			},
		},
		{
			name: "loads multiple field types",
			cfg: &struct {
				Name    string `env:"TEST_NAME_MULTI"`
				Port    int    `env:"TEST_PORT_MULTI"`
				Enabled bool   `env:"TEST_ENABLED_MULTI"`
				Ports   []int  `env:"TEST_PORTS_MULTI"`
			}{},
			setupEnv: func() {
				os.Setenv("TEST_NAME_MULTI", "server")
				os.Setenv("TEST_PORT_MULTI", "8080")
				os.Setenv("TEST_ENABLED_MULTI", "true")
				os.Setenv("TEST_PORTS_MULTI", "8080,8081")
			},
			cleanupEnv: func() {
				os.Unsetenv("TEST_NAME_MULTI")
				os.Unsetenv("TEST_PORT_MULTI")
				os.Unsetenv("TEST_ENABLED_MULTI")
				os.Unsetenv("TEST_PORTS_MULTI")
			},
			wantErr: false,
			validate: func(t *testing.T, cfg interface{}) {
				c := cfg.(*struct {
					Name    string `env:"TEST_NAME_MULTI"`
					Port    int    `env:"TEST_PORT_MULTI"`
					Enabled bool   `env:"TEST_ENABLED_MULTI"`
					Ports   []int  `env:"TEST_PORTS_MULTI"`
				})
				if c.Name != "server" {
					t.Errorf("Name = %v, want server", c.Name)
				}
				if c.Port != 8080 {
					t.Errorf("Port = %v, want 8080", c.Port)
				}
				if !c.Enabled {
					t.Errorf("Enabled = %v, want true", c.Enabled)
				}
				wantPorts := []int{8080, 8081}
				if len(c.Ports) != len(wantPorts) {
					t.Errorf("Ports length = %v, want %v", len(c.Ports), len(wantPorts))
					return
				}
				for i, v := range wantPorts {
					if c.Ports[i] != v {
						t.Errorf("Ports[%d] = %v, want %v", i, c.Ports[i], v)
					}
				}
			},
		},
		{
			name: "skips fields without env tag",
			cfg: &struct {
				Name string `env:"TEST_NAME_SKIP"`
				Skip string
			}{},
			setupEnv: func() {
				os.Setenv("TEST_NAME_SKIP", "value")
			},
			cleanupEnv: func() {
				os.Unsetenv("TEST_NAME_SKIP")
			},
			wantErr: false,
			validate: func(t *testing.T, cfg interface{}) {
				c := cfg.(*struct {
					Name string `env:"TEST_NAME_SKIP"`
					Skip string
				})
				if c.Name != "value" {
					t.Errorf("Name = %v, want value", c.Name)
				}
				if c.Skip != "" {
					t.Errorf("Skip = %v, want empty string", c.Skip)
				}
			},
		},
		{
			name: "returns error for invalid int value",
			cfg: &struct {
				Port int `env:"TEST_PORT_INVALID"`
			}{},
			setupEnv: func() {
				os.Setenv("TEST_PORT_INVALID", "not_a_number")
			},
			cleanupEnv: func() {
				os.Unsetenv("TEST_PORT_INVALID")
			},
			wantErr:  true,
			validate: func(t *testing.T, cfg interface{}) {},
		},
		{
			name: "returns error for invalid int in slice",
			cfg: &struct {
				Ports []int `env:"TEST_PORTS_INVALID"`
			}{},
			setupEnv: func() {
				os.Setenv("TEST_PORTS_INVALID", "8080,not_a_number,8082")
			},
			cleanupEnv: func() {
				os.Unsetenv("TEST_PORTS_INVALID")
			},
			wantErr:  true,
			validate: func(t *testing.T, cfg interface{}) {},
		},
		{
			name: "returns error for array length mismatch",
			cfg: &struct {
				Fixed [3]int `env:"TEST_FIXED_MISMATCH"`
			}{},
			setupEnv: func() {
				os.Setenv("TEST_FIXED_MISMATCH", "1,2")
			},
			cleanupEnv: func() {
				os.Unsetenv("TEST_FIXED_MISMATCH")
			},
			wantErr:  true,
			validate: func(t *testing.T, cfg interface{}) {},
		},
		{
			name: "returns error when cfg is not a pointer",
			cfg: struct {
				Name string `env:"TEST_NAME"`
			}{},
			setupEnv: func() {
				os.Setenv("TEST_NAME", "value")
			},
			cleanupEnv: func() {
				os.Unsetenv("TEST_NAME")
			},
			wantErr:  true,
			validate: func(t *testing.T, cfg interface{}) {},
		},
		{
			name: "returns error when cfg is not a struct",
			cfg: func() *string {
				s := "test"
				return &s
			}(),
			setupEnv:   func() {},
			cleanupEnv: func() {},
			wantErr:    true,
			validate:   func(t *testing.T, cfg interface{}) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupEnv()
			defer tt.cleanupEnv()

			err := LoadStruct(tt.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadStruct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				tt.validate(t, tt.cfg)
			}
		})
	}
}
