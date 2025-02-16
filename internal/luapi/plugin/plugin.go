package plugin

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/gvcgo/goutils/pkgs/gutils"
	"github.com/gvcgo/version-manager/internal/cnf"
	"github.com/gvcgo/version-manager/internal/luapi/lua_global"
	"github.com/gvcgo/version-manager/internal/tui/table"
	"github.com/gvcgo/version-manager/internal/utils"
)

type Plugin struct {
	FileName      string `json:"file_name"`
	PluginName    string `json:"plugin_name"`
	PluginVersion string `json:"plugin_version"`
	SDKName       string `json:"sdk_name"`
	Prequisite    string `json:"prequisite"`
	Homepage      string `json:"homepage"`
}

type Plugins struct {
	pls map[string]Plugin
}

func NewPlugins() *Plugins {
	p := &Plugins{
		pls: make(map[string]Plugin),
	}
	if ok, _ := gutils.PathIsExist(cnf.GetPluginDir()); !ok {
		p.Update()
	}
	return p
}

func (p *Plugins) Update() {
	UpdatePlugins()
}

func (p *Plugins) LoadAll() {
	pDir := cnf.GetPluginDir()
	files, _ := os.ReadDir(pDir)

	for _, f := range files {
		if f.IsDir() || !strings.HasSuffix(f.Name(), ".lua") {
			continue
		}
		pl := Plugin{
			FileName: f.Name(),
		}
		ll := lua_global.NewLua()
		L := ll.L
		L.DoFile(filepath.Join(pDir, f.Name()))
		pl.PluginName = GetConfItemFromLua(L, PluginName)
		if pl.PluginName == "" {
			continue
		}
		pl.PluginVersion = GetConfItemFromLua(L, PluginVersion)
		pl.SDKName = GetConfItemFromLua(L, SDKName)
		if pl.SDKName == "" {
			continue
		}
		pl.Prequisite = GetConfItemFromLua(L, Prequisite)
		pl.Homepage = GetConfItemFromLua(L, Homepage)
		if pl.Homepage == "" {
			continue
		}
		if !DoLuaItemExist(L, InstallerConfig) || !DoLuaItemExist(L, Crawler) {
			continue
		}
		p.pls[pl.PluginName] = pl
		ll.Close()
	}
}

/*
TODO:
1. run a lua plugin
2. cache version list
3. show plugin list
*/
func (p *Plugins) GetPlugin(pName string) Plugin {
	p.LoadAll()
	if pl, ok := p.pls[pName]; ok {
		return pl
	}
	return Plugin{}
}

func (p *Plugins) GetPluginList() (pl []Plugin) {
	p.LoadAll()
	for _, v := range p.pls {
		pl = append(pl, v)
	}
	return
}

func (p *Plugins) GetPluginSortedRows() (rows []table.Row) {
	p.LoadAll()
	for _, v := range p.pls {
		rows = append(rows, table.Row{
			v.PluginName,
			v.Homepage,
			v.PluginVersion,
		})
	}
	utils.SortVersionAscend(rows)
	return
}
