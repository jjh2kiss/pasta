package sched

import "testing"
import "strings"
import "reflect"
import "errors"

func TestNewInitProcessInternal(t *testing.T) {
	testcases := []struct {
		title                 string
		lines                 string
		expected_init_process *InitProcess
		expected_error        error
	}{

		{
			"Single line, valid",
			"systemd (1, #threads: 1)",
			&InitProcess{"systemd", 1, 1},
			nil,
		},
		{
			"Multiline, Valid",
			`systemd (1, #threads: 1)
			-------------------------------------------------------------------
			se.exec_start                                :     958366708.536261
			se.vruntime                                  :          7841.031737
			se.sum_exec_runtime                          :         18505.036860
			se.statistics.sum_sleep_runtime              :     958347872.252774
			`,
			&InitProcess{"systemd", 1, 1},
			nil,
		},
		{
			"Single line, valid, init",
			"init (1, #threads: 1)",
			&InitProcess{"init", 1, 1},
			nil,
		},
		{
			"Single line, valid, init, pid 1234, thread 1357",
			"init (1234, #threads: 1357)",
			&InitProcess{"init", 1234, 1357},
			nil,
		},
		{
			"Invalid, empty firstline",
			"",
			nil,
			errors.New("Can't read data from reader"),
		},
		{
			"Invalid, wrong firstline",
			"init",
			nil,
			errors.New("fail to parse first_line(init)"),
		},
		{
			"Invalid, without threads field",
			"init (1)",
			nil,
			errors.New("fail to parse other fields(1)"),
		},
		{
			"Invalid, string pid",
			"init (a, #threads: 10)",
			nil,
			errors.New("strconv.ParseInt: parsing \"a\": invalid syntax"),
		},
		{
			"Invalid, threads format",
			"init (1, #threads 10)",
			nil,
			errors.New("fail to parse threads field(#threads 10)"),
		},
		{
			"Invalid, string threads",
			"init (1, #threads: a)",
			nil,
			errors.New("strconv.ParseInt: parsing \"a\": invalid syntax"),
		},
	}

	for _, testcase := range testcases {
		reader := strings.NewReader(testcase.lines)
		actual_init_process, actual_error := newInitProcess(reader)

		if reflect.DeepEqual(actual_init_process, testcase.expected_init_process) == false {
			t.Errorf("%s\nexpected \n%v\nBut\n%v\n",
				testcase.title,
				testcase.expected_init_process,
				actual_init_process,
			)
			return
		}

		if actual_error != nil && testcase.expected_error != nil {
			if actual_error.Error() != testcase.expected_error.Error() {
				t.Errorf("%s\nexpected error \n%v\nBut\n%v\n",
					testcase.title,
					testcase.expected_error,
					actual_error,
				)
			}
		} else if actual_error == nil && testcase.expected_error == nil {
			continue
		} else {
			t.Errorf("%s\nexpected error is\n%v\nBut\n%v\n",
				testcase.title,
				testcase.expected_error,
				actual_error,
			)
		}

	}
}

func TestGetInitProcess(t *testing.T) {
	init_process, err := GetInitProcess()
	if err != nil {
		t.Error(err.Error())
	}

	if len(init_process.Name) == 0 {
		t.Errorf("length of name of init process should be greater than 0, but 0")
	}

	if init_process.Pid == 0 {
		t.Errorf("pid of init process should be greater than 0, but 0")
	}

}

func TestRunningInContainer(t *testing.T) {
	testcases := []struct {
		in       *InitProcess
		expected bool
	}{
		{&InitProcess{"init", 1, 1}, false},
		{&InitProcess{"init", 2, 1}, true},
		{&InitProcess{"init", 1234, 1}, true},
	}

	for _, testcase := range testcases {
		actual := testcase.in.RunningInContainer()
		if actual != testcase.expected {
			t.Errorf("expected %v, but %v\n", testcase.expected, actual)
		}
	}

}
