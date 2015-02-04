========================
REST API Documentation
========================

CPU
---

CPU Info
++++++++

Get CPU information.

.. code-block:: bash

    GET /cpuinfo HTTP/1.1

Memory
------

Memory Info
+++++++++++

Get memory information.

.. code-block:: bash

    GET /memstats HTTP/1.1

Battery
-------

Battery Info
++++++++++++

Get battery information.

.. code-block:: bash

    GET /battery HTTP/1.1

Uptime
------

System Uptime Info
++++++++++++++++++

.. code-block:: bash

    GET /sysinfo HTTP/1.1

Sysfs
-----

Sysfs File System
+++++++++++++++++

.. code-block:: bash

    GET /sysfs HTTP/1.1

Proceess
--------

Proceess Memory Info
++++++++++++++++++++

Lists all the processes, executables, and shared libraries that are using up virtual memory.

.. code-block:: bash

    GET /meminfo HTTP/1.1

Network
-------

Interface Info
++++++++++++++

Get network interface information.

.. code-block:: bash

    GET /ifconfig HTTP/1.1

ARP Table
+++++++++

.. code-block:: bash

    GET /arp HTTP/1.1

Routing Table
+++++++++++++

.. code-block:: bash

    GET /route HTTP/1.1
