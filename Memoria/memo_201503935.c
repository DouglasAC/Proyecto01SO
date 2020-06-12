#include <linux/kernel.h>   
#include <linux/module.h>
#include <linux/proc_fs.h>
#include <linux/seq_file.h> 
#include <asm/uaccess.h> 
#include <linux/hugetlb.h>
#include <linux/init.h>
#include <linux/fs.h>

MODULE_LICENSE("GPL");
MODULE_AUTHOR("201503935 - Douglas Daniel Aguilar Cuque");
MODULE_DESCRIPTION("Obtener y escribir el estado de la memoria RAM");

#define ENTRY_NAME "memo_201503935"
#define PERMS 0644 
#define PARENT NULL 
#define BUFSIZE  	150

struct sysinfo info_ram;

static int escribirArchivo(struct seq_file * archivo, void *v) {	
    long memoriaT, memoriaL;
    si_meminfo(&info_ram);
    memoriaT = info_ram.totalram * 4;
    memoriaL = info_ram.freeram * 4;
    seq_printf(archivo, "Carnet: 201503935\n");
    seq_printf(archivo, "Nombre: Douglas Daniel Aguilar Cuque\n");
    seq_printf(archivo, "Memoria Total: %1lu MB\n", memoriaT / 1024);
    seq_printf(archivo, "Memoria Libre: %2lu MB \n",  memoriaL / 1024);
    seq_printf(archivo, "Memoria en uso: %li %%\n", (memoriaL * 100)/memoriaT) ;
    return 0;
}

static int abrir(struct inode *inode, struct  file *file) {
  return single_open(file, escribirArchivo, NULL);
}

static struct file_operations ope =
{    
    .open = abrir,
    .read = seq_read
};

static int inicio(void)
{
    printk(KERN_INFO "Carnet: 201503935\n");
    if(!proc_create(ENTRY_NAME, PERMS, NULL, &ope))
    {
        printk("Error al crear ");
        remove_proc_entry(ENTRY_NAME,NULL);
        return -ENOMEM;
    }
    return 0;
}
 
static void fin(void)
{
    remove_proc_entry(ENTRY_NAME, NULL);
    printk(KERN_INFO "Curso: Sistemas Operativos 1\n");
}
 
module_init(inicio);
module_exit(fin); 