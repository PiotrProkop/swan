import glog as log
import subprocess
import time
import copy


class Shell:
    """
    Convenience class to run a number of commands simultaneously and reap their exit statuses.
    """
    def __init__(self, commands, output="output.txt"):
        """
        :param commands: List of commands to run
        :param output: File to save commands output in.
        """
        self.processes = {}

        for command in commands:
            if command == "":
                log.warning("Command in list %s is empty: aborting execution" % str(commands))
                return

            output_file = open(output, "a+")
            log.info("started command: '" + str(command) + "'")
            p = subprocess.Popen(["sh", "-c", str(command)], stdout=output_file, stderr=output_file)
            self.processes[p.pid] = {'process': p, 'command': command, 'output_file': output_file, 'status': None}

        # Make a copy of processes. Otherwise, we loose hold of the process objects when removing from the running
        # list.
        running = copy.copy(self.processes)

        # Reap process statuses.
        while len(running) is not 0:
            exited_pids = []
            for pid, process_obj in running.iteritems():
                process = process_obj['process']
                status = process.poll()
                if status is not None:
                    command = process_obj['command']
                    log.info("ended command: '" + str(command) + "' with status code " + str(status))
                    output_file = process_obj['output_file']
                    output_file.flush()
                    exited_pids.append(pid)

                    # Update original record
                    self.processes[pid]['status'] = status

                    # TODO: Write post mortem to log

            for pid in exited_pids:
                del running[pid]

            # Prevent busy loop
            time.sleep(0.1)


class Delay:
    """
    Convenience class to delay execution of a command with a configurable number of seconds.
    """
    def __init__(self, seconds, command):
        """
        :param seconds: Period in seconds to wait before executing
        :param command: Command to run
        """
        self.seconds = seconds
        self.command = command

    def __str__(self):
        if self.seconds <= 0:
            return self.command

        return "sleep %d && %s" % (self.seconds, self.command)


class RunFor:
    """
    Convenience class to run a command for a specified number of seconds and send a configurable signal.
    """
    def __init__(self, seconds, command, signal="SIGINT"):
        """
        :param seconds: Period in seconds to execute
        :param command: Command to run
        :param signal: Signal to send. For example SIGINT, SIGHUP, etc.
        """
        self.seconds = seconds
        self.command = command
        self.signal = signal

    def __str__(self):
        return "timeout -s %s %d %s" % (self.signal, self.seconds, self.command)
