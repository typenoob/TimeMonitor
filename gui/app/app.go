package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"gui/utils"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type Singleton struct {
	a                       fyne.App
	w                       fyne.Window
	c                       *fyne.Container
	main_manu               *fyne.MainMenu
	menu_item_quit          *fyne.MenuItem
	menu_item_setting       *fyne.MenuItem
	container_setting       *fyne.Container
	container_refresh       *fyne.Container
	container_search        *fyne.Container
	wg_table                *widget.Table
	wg_popup_setting        *widget.PopUp
	wg_entry_ntp_server     *widget.Entry
	wg_entry_interval       *widget.Entry
	wg_entry_policy         *widget.Entry
	wg_entry_search         *widget.Entry
	wg_label_ntp_server     *widget.Label
	wg_label_refresh        *widget.Label
	wg_label_refresh_data   *widget.Label
	wg_label_interval       *widget.Label
	wg_label_policy         *widget.Label
	wg_label_policy_postfix *widget.Label
	wg_button_save          *widget.Button
	wg_button_refresh       *widget.Button
	wg_button_search        *widget.Button
	rs_refresh              fyne.Resource
	rs_search               fyne.Resource
	rs_green_dot            fyne.Resource
	rs_red_dot              fyne.Resource
	last_refresh_time       binding.String
	t                       *time.Ticker
	filter_key              string
	records                 []utils.Record
}

var (
	instance *Singleton
	once     sync.Once
)

const (
	APP_KEY                  = "CLOCK"
	ELECTRONIC_CLOCK_MONITOR = "电子时钟监控系统"
	NTP_SERVER_PLACEHOLDER   = "请输入NTP服务器地址"
	NTP_SERVER_INVALID       = "NTP服务器地址不合法"
	NTP_SERVER_CONN_TIMEOUT  = "NTP服务器连接超时"
	NTP_SERVER_ADDRESS       = "NTP服务器地址"
	LAST_REFRESH_TIME        = "上次刷新时间"
	REFRESH_INTERVAL         = "刷新时间间隔(秒)"
	INTERVAL_PLACEHOLDER     = "请输入刷新时间间隔(秒)"
	INTERVAL_INVALID         = "时间间隔不合法"
	DEVICE_ACTIVE_POLICY     = "设备在线策略(小时)"
	DEVICE_ACTIVE_POSTFIX    = "小时内成功更新"
	SEARCH_PLACEHOLDER       = "请输入要搜索的内容"
	ID                       = "编号"
	IP_ADDRESS               = "IP地址"
	LAST_OK_TIME             = "最后一次成功更新时间"
	STATE                    = "状态"
	SYSTEM                   = "系统"
	OPTION                   = "选项"
	SETTING                  = "设置"
	EXIT                     = "退出"
	SAVE                     = "保存"
	REFRESH_ICON_PATH        = "assets/refresh.svg"
	SEARCH_ICON_PATH         = "assets/search.svg"
	GREEN_DOT_ICON_PATH      = "assets/green-dot.svg"
	RED_DOT_ICON_PATH        = "assets/red-dot.svg"
)

func (s *Singleton) SetNTPServer(value string) {
	s.a.Preferences().SetString("NTPServer", value)
}

func (s *Singleton) GetNTPServer() string {
	return s.a.Preferences().String("NTPServer")
}

func (s *Singleton) SetInterval(value int) {
	s.a.Preferences().SetInt("Interval", value)
	s.t.Reset(time.Duration(s.GetInterval()) * time.Second)
}

func (s *Singleton) GetInterval() int {
	return s.a.Preferences().IntWithFallback("Interval", 60)
}

func (s *Singleton) SetPolicy(value int) {
	s.a.Preferences().SetInt("Policy", value)
}

func (s *Singleton) GetPolicy() int {
	return s.a.Preferences().IntWithFallback("Policy", 2)
}
func (s *Singleton) SetLastRefreshTime(t time.Time) {
	s.last_refresh_time.Set(t.Format(time.DateTime))
}

func (s *Singleton) GetState(t time.Time) bool {
	return time.Since(t) < time.Hour*time.Duration(s.GetPolicy())
}

func (s *Singleton) FetchRecords() {
	if records, err := utils.GetAllRecord(s.GetNTPServer()); err != nil {
	} else {
		if j, err := json.Marshal(records); err != nil {
			log.Fatal(err)
		} else {
			s.a.Preferences().SetString("Records", string(j))
		}
	}
}

func (s *Singleton) OnRefresh() {
	s.FetchRecords()
	s.records = s.GetRecords()
	s.SetLastRefreshTime(time.Now())
	s.wg_table.Refresh()
}

func (s *Singleton) OnSearch(value string) {
	s.filter_key = value
	s.wg_entry_search.SetText("")
	s.records = s.GetRecords()
	s.wg_table.Refresh()
}

func (s *Singleton) GetRecords() []utils.Record {
	var records []utils.Record
	json.Unmarshal([]byte(s.a.Preferences().String("Records")), &records)
	records_filtered := make([]utils.Record, 0)
	for _, item := range records {
		if strings.Contains(item.ID, s.filter_key) || strings.Contains(item.IPAddress, s.filter_key) || strings.Contains(item.LastOkTime, s.filter_key) {
			records_filtered = append(records_filtered, item)
		}
	}
	return records_filtered
}

func NTPServerValidator(s string) error {
	if net.ParseIP(s) == nil {
		if !utils.IsValidDomainName(s) {
			return errors.New(NTP_SERVER_INVALID)
		}
	}
	return nil
}

func NumberValidator(s string) error {
	if _, err := strconv.Atoi(s); err != nil {
		return errors.New(INTERVAL_INVALID)
	}
	return nil
}

func (s *Singleton) OnSettingTap() {
	s.wg_entry_ntp_server.SetText(s.a.Preferences().String("NTPServer"))
	s.wg_entry_ntp_server.SetPlaceHolder(NTP_SERVER_PLACEHOLDER)
	s.wg_entry_ntp_server.Validator = NTPServerValidator
	s.wg_entry_interval.SetText(strconv.Itoa(s.GetInterval()))
	s.wg_entry_interval.SetPlaceHolder(INTERVAL_PLACEHOLDER)
	s.wg_entry_interval.Validator = NumberValidator
	s.wg_entry_policy.SetText(strconv.Itoa(s.GetPolicy()))
	s.wg_entry_interval.Validator = NumberValidator
	s.wg_popup_setting.Resize(fyne.NewSize(330, 440))
	s.wg_popup_setting.Show()
}

func (s *Singleton) OnSave() {
	if s.wg_entry_ntp_server.Text == s.GetNTPServer() && s.wg_entry_interval.Text == strconv.Itoa(s.GetInterval()) {
		s.wg_popup_setting.Hide()
		return
	}
	if s.wg_entry_ntp_server.Validate() != nil {
		dialog.NewError(s.wg_entry_ntp_server.Validate(), s.w).Show()
		return
	} else {
		if _, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", s.wg_entry_ntp_server.Text, utils.Port), 1*time.Second); err != nil {
			dialog.NewError(errors.New(NTP_SERVER_CONN_TIMEOUT), s.w).Show()
			return
		} else {
			s.SetNTPServer(s.wg_entry_ntp_server.Text)
		}
	}
	if s.wg_entry_interval.Validate() != nil {
		dialog.NewError(s.wg_entry_interval.Validate(), s.w).Show()
		return
	} else if i, err := strconv.Atoi(s.wg_entry_interval.Text); err == nil {
		s.SetInterval(i)
	}
	if s.wg_entry_policy.Validate() != nil {
		dialog.NewError(s.wg_entry_policy.Validate(), s.w).Show()
		return
	} else if i, err := strconv.Atoi(s.wg_entry_policy.Text); err == nil {
		s.SetPolicy(i)
	}
	s.wg_popup_setting.Hide()
	s.OnRefresh()
}

func (s *Singleton) Start() {
	s.t = time.NewTicker(time.Duration(s.GetInterval()) * time.Second)
	s.OnRefresh()
	go func() {
		for {
			<-s.t.C
			s.OnRefresh()
		}
	}()
	s.w.ShowAndRun()
}

func (s *Singleton) OnClose() {
	s.w.Close()
	os.Exit(0)
}

func (s *Singleton) initRefresh() {
	s.last_refresh_time = binding.BindPreferenceString("LastRefreshTime", s.a.Preferences())
	s.wg_label_refresh = widget.NewLabel(LAST_REFRESH_TIME)
	s.wg_label_refresh_data = widget.NewLabelWithData(s.last_refresh_time)
	s.wg_button_refresh = widget.NewButtonWithIcon("", s.rs_refresh, s.OnRefresh)
	s.container_refresh = container.NewHBox(container.New(layout.NewCustomPaddedHBoxLayout(theme.Padding()-theme.InnerPadding()), s.wg_label_refresh, s.wg_label_refresh_data), s.wg_button_refresh)
}

func (s *Singleton) initSearch() {
	s.wg_entry_search = widget.NewEntry()
	s.wg_entry_search.SetPlaceHolder(SEARCH_PLACEHOLDER)
	s.wg_button_search = widget.NewButtonWithIcon("", s.rs_search, func() {
		s.OnSearch(s.wg_entry_search.Text)
	})
	s.container_search = container.NewBorder(nil, nil, nil, s.wg_button_search, container.New(layout.NewCustomPaddedLayout(0, 0, theme.Padding()*8, 0), s.wg_entry_search))
}

func (s *Singleton) initTable() {
	var title = []string{ID, IP_ADDRESS, LAST_OK_TIME, STATE}
	s.records = s.GetRecords()
	s.wg_table = widget.NewTable(
		func() (int, int) { return len(s.records), 4 },
		func() fyne.CanvasObject { return container.NewStack(widget.NewLabel(SYSTEM), widget.NewIcon(nil)) },
		func(i widget.TableCellID, o fyne.CanvasObject) {
			label := o.(*fyne.Container).Objects[0].(*widget.Label)
			icon := o.(*fyne.Container).Objects[1].(*widget.Icon)
			label.Show()
			icon.Hide()
			switch i.Col {
			case 0:
				label.SetText(s.records[i.Row].ID)
			case 1:
				label.SetText(s.records[i.Row].IPAddress)
			case 2:
				if t, err := time.Parse(time.RFC3339, s.records[i.Row].LastOkTime); err != nil {
					log.Fatal(err)
				} else {
					label.SetText(t.Format(time.DateTime))
				}
			case 3:
				label.Hide()
				icon.Show()
				if t, err := time.Parse(time.RFC3339, s.records[i.Row].LastOkTime); err != nil {
					log.Fatal(err)
				} else {
					if s.GetState(t) {
						icon.SetResource(s.rs_green_dot)
					} else {
						icon.SetResource(s.rs_red_dot)
					}
				}
			}
		})
	s.wg_table.CreateHeader = func() fyne.CanvasObject {
		return widget.NewLabel("")
	}
	s.wg_table.UpdateHeader = func(i widget.TableCellID, o fyne.CanvasObject) {
		o.(*widget.Label).SetText(title[i.Col])
	}
	s.wg_table.ShowHeaderRow = true
	s.wg_table.SetColumnWidth(1, 100)
	s.wg_table.SetColumnWidth(2, 200)
}

func (s *Singleton) initMainMenu() {
	s.menu_item_quit = fyne.NewMenuItem(EXIT, s.OnClose)
	s.menu_item_quit.IsQuit = true
	s.menu_item_setting = fyne.NewMenuItem(SETTING, s.OnSettingTap)
	s.main_manu = fyne.NewMainMenu(fyne.NewMenu(OPTION, s.menu_item_setting, s.menu_item_quit))
	s.w.SetMainMenu(s.main_manu)
}

func (s *Singleton) initPopUpSetting() {
	s.wg_label_ntp_server = widget.NewLabel(NTP_SERVER_ADDRESS)
	s.wg_entry_ntp_server = widget.NewEntry()
	s.wg_label_interval = widget.NewLabel(REFRESH_INTERVAL)
	s.wg_entry_interval = widget.NewEntry()
	s.wg_label_policy = widget.NewLabel(DEVICE_ACTIVE_POLICY)
	s.wg_entry_policy = widget.NewEntry()
	s.wg_label_policy_postfix = widget.NewLabel(DEVICE_ACTIVE_POSTFIX)
	s.wg_button_save = widget.NewButton(SAVE, s.OnSave)
	s.container_setting = container.NewVBox(container.New(layout.NewFormLayout(), s.wg_label_ntp_server, s.wg_entry_ntp_server, s.wg_label_interval, s.wg_entry_interval, s.wg_label_policy, container.NewBorder(nil, nil, nil, s.wg_label_policy_postfix, s.wg_entry_policy)), layout.NewSpacer(), s.wg_button_save)
	s.wg_popup_setting = widget.NewPopUp(s.container_setting, s.w.Canvas())
}

func (s *Singleton) initContainer() {
	s.c = container.NewBorder(container.NewGridWithColumns(2, s.container_refresh, s.container_search), layout.NewSpacer(), layout.NewSpacer(), layout.NewSpacer(), instance.wg_table)
}

func (s *Singleton) loadResource() {
	s.rs_refresh = theme.Icon(theme.IconNameViewRefresh)
	s.rs_search = theme.Icon(theme.IconNameSearch)
	if rs, err := fyne.LoadResourceFromPath(GREEN_DOT_ICON_PATH); err != nil {
		s.rs_green_dot = theme.Icon(theme.IconNameBrokenImage)
	} else {
		s.rs_green_dot = rs
	}
	if rs, err := fyne.LoadResourceFromPath(RED_DOT_ICON_PATH); err != nil {
		s.rs_red_dot = theme.Icon(theme.IconNameBrokenImage)
	} else {
		s.rs_red_dot = rs
	}
}

func GetInstance() *Singleton {
	once.Do(func() {
		a := app.NewWithID(APP_KEY)
		w := a.NewWindow(ELECTRONIC_CLOCK_MONITOR)
		instance = &Singleton{a: a, w: w}
		instance.loadResource()
		instance.initMainMenu()
		instance.initRefresh()
		instance.initSearch()
		instance.initTable()
		instance.initPopUpSetting()
		instance.initContainer()
		instance.w.SetContent(instance.c)
		instance.w.Resize(fyne.NewSize(650, 950))
	})
	return instance
}
